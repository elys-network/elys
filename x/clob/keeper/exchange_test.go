package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (suite *KeeperTestSuite) TestExchange() {
	suite.ResetSuite()
	baseDenom := "uatom"

	marketId := uint64(1)
	balancePerAccount := sdk.NewCoins(sdk.NewInt64Coin("uusdc", 1000_000_0))
	subAccounts := suite.SetupSubAccounts(2, balancePerAccount)

	buyerSubAccount := subAccounts[0]
	sellerSubAccount := subAccounts[1]

	t1 := types.Trade{
		BuyerSubAccount:  buyerSubAccount,
		SellerSubAccount: sellerSubAccount,
		MarketId:         1,
		Price:            math.LegacyNewDecWithPrec(101, 1),
		Quantity:         math.NewInt(100),
	}
	testCases := []struct {
		name           string
		expectedErrMsg string
		input          types.Trade
		pre            func()
		post           func()
	}{
		{
			"no market",
			types.ErrPerpetualMarketNotFound.Error(),
			t1,
			func() {
			},
			func() {
			},
		},
		{
			"buyer lack funds",
			"insufficient funds",
			t1,
			func() {
				suite.CreateMarket(baseDenom)
				err := suite.app.ClobKeeper.SendFromSubAccount(suite.ctx, buyerSubAccount, authtypes.NewModuleAddress(types.ModuleName), balancePerAccount)
				suite.Require().NoError(err)
			},
			func() {
			},
		},
		{
			"both subAccounts never opened a position before",
			"",
			t1,
			func() {
				err := suite.app.ClobKeeper.AddToSubAccount(suite.ctx, authtypes.NewModuleAddress(types.ModuleName), buyerSubAccount, balancePerAccount)
				suite.Require().NoError(err)
			},
			func() {
				lastAvgTradePrice := suite.app.ClobKeeper.GetLastAverageTradePrice(suite.ctx, 1)
				suite.Require().Equal(math.LegacyNewDecWithPrec(1010, 2), lastAvgTradePrice)

				market, err := suite.app.ClobKeeper.GetPerpetualMarket(suite.ctx, marketId)
				suite.Require().NoError(err)

				// new quantity gets created
				suite.Require().Equal(t1.Quantity, market.TotalOpen)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.IncreaseHeight(1)
			tc.pre()
			err := suite.app.ClobKeeper.Exchange(suite.ctx, tc.input)
			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
			}
			tc.post()
		})
	}
}
