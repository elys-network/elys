package keeper_test

import (
	"github.com/elys-network/elys/testutil/sample"
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (suite *KeeperTestSuite) TestDebt() {
	for _, tc := range []struct {
		desc              string
		senderInitBalance sdk.Coins
		moduleInitBalance sdk.Coins
		unbondAmount      math.Int
		expSenderBalance  sdk.Coins
		expPass           bool
	}{
		{
			desc:              "successful debt process",
			senderInitBalance: sdk.Coins{sdk.NewInt64Coin(types.GetShareDenom(), 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100000000)},
			moduleInitBalance: sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)},
			unbondAmount:      math.NewInt(1000000),
			expSenderBalance:  sdk.Coins{sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)}.Sort(),
			expPass:           true,
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			// bootstrap accounts
			sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

			// bootstrap balances
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.moduleInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, minttypes.ModuleName, types.ModuleName, tc.moduleInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.senderInitBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, tc.senderInitBalance)
			suite.Require().NoError(err)

			shareDenom := types.GetShareDenom()

			// Set an entity to assetprofile
			entry := assetprofiletypes.Entry{
				Authority:       authtypes.NewModuleAddress(types.ModuleName).String(),
				BaseDenom:       shareDenom,
				Decimals:        ptypes.BASE_DECIMAL,
				Denom:           shareDenom,
				DisplayName:     shareDenom,
				CommitEnabled:   true,
				WithdrawEnabled: true,
			}
			suite.app.AssetprofileKeeper.SetEntry(suite.ctx, entry)

			// Commit LP token
			err = suite.app.CommitmentKeeper.CommitLiquidTokens(
				suite.ctx, sender,
				tc.senderInitBalance[0].Denom, tc.senderInitBalance[0].Amount,
				uint64(suite.ctx.BlockTime().Unix()),
			)
			suite.Require().NoError(err)

			params := suite.app.StablestakeKeeper.GetParams(suite.ctx)
			params.TotalValue = math.NewInt(10)
			params.InterestRate = math.LegacyNewDec(10)
			suite.app.StablestakeKeeper.SetParams(suite.ctx, params)

			err = suite.app.StablestakeKeeper.Borrow(suite.ctx, sender, sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)), 1)
			suite.Require().NoError(err)
			suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Hour * 24 * 365)).WithBlockHeight(10)

			// Pay partial
			err = suite.app.StablestakeKeeper.Repay(suite.ctx, sender, sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10)), 1)
			suite.Require().NoError(err)

			res := suite.app.StablestakeKeeper.UpdateInterestAndGetDebt(suite.ctx, sender, 1, ptypes.BaseCurrency)
			suite.Require().Equal(res.Borrowed.String(), "1000")
			suite.Require().Equal(res.InterestStacked.String(), "10000")
			suite.Require().Equal(res.InterestPaid.String(), "10")
			allDebts := suite.app.StablestakeKeeper.GetAllDebts(suite.ctx)
			suite.Require().Len(allDebts, 1)

			// Pay rest, ensure we don't pay multiple times
			err = suite.app.StablestakeKeeper.Repay(suite.ctx, sender, sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10990)), 1)
			suite.Require().NoError(err)
			res = suite.app.StablestakeKeeper.GetDebt(suite.ctx, sender)
			suite.Require().Equal(res.Borrowed.String(), "0")
			suite.Require().Equal(res.InterestStacked.String(), "0")
			suite.Require().Equal(res.InterestPaid.String(), "0")
		})
	}
}

func (suite *KeeperTestSuite) TestCloseOnUnableToRepay() {
	p := types.AmmPool{
		Id:               1,
		TotalLiabilities: sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)},
	}

	suite.app.StablestakeKeeper.SetAmmPool(suite.ctx, p)

	debt := types.Debt{
		Address:               sample.AccAddress(),
		Borrowed:              math.NewInt(100),
		InterestPaid:          math.NewInt(10),
		InterestStacked:       math.NewInt(50),
		BorrowTime:            1,
		LastInterestCalcTime:  1,
		LastInterestCalcBlock: 1,
	}
	suite.app.StablestakeKeeper.SetDebt(suite.ctx, debt)
	suite.app.StablestakeKeeper.CloseOnUnableToRepay(suite.ctx, debt.GetOwnerAccount(), 1, sdk.DefaultBondDenom)

	r := suite.app.StablestakeKeeper.GetAmmPool(suite.ctx, 1)
	suite.Assert().Equal(types.AmmPool{
		Id:               1,
		TotalLiabilities: sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 860)},
	}, r)

}
