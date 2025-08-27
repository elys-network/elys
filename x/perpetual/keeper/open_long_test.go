package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/v7/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/v7/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *PerpetualKeeperTestSuite) TestOpenLong() {
	suite.SetupCoinPrices()
	addr := suite.AddAccounts(10, nil)
	amount := math.NewInt(1000)
	poolCreator := addr[0]
	positionCreator := addr[1]
	poolId := uint64(1)
	var ammPool ammtypes.Pool
	msg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(2),
		Position:        types.Position_LONG,
		PoolId:          poolId,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyOneDec().MulInt64(6),
		StopLossPrice:   math.LegacyZeroDec(),
	}
	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{
			"amm pool not found",
			"pool does not exist",
			func() {
				ammPool = suite.CreateNewAmmPool(poolCreator, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
				poolId = ammPool.PoolId
				enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Pool: leveragelpmoduletypes.AddPool{
						AmmPoolId:       poolId,
						LeverageMax:     math.LegacyMustNewDecFromStr("10"),
						AdlTriggerRatio: math.LegacyNewDec(1),
					},
				}
				_, err := leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
				suite.Require().NoError(err)

				suite.app.AmmKeeper.RemovePool(suite.ctx, ammPool.PoolId)
			},
		},
		{
			"collateral amount is too high",
			"borrowed amount is higher than pool depth",
			func() {
				suite.app.AmmKeeper.SetPool(suite.ctx, ammPool)
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.BorrowInterestRateMin = math.LegacyZeroDec()
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = msg.Collateral.Amount.MulRaw(1000_000_000)
			},
		},
		{
			"Borrow fails: lack of funds",
			"user does not have enough balance of the required coin",
			func() {
				msg.Collateral.Amount = msg.Collateral.Amount.QuoRaw(1000_000_000)
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.BorrowInterestRateMin = math.LegacyMustNewDecFromStr("0.12")
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, positionCreator, govtypes.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, suite.GetAccountIssueAmount())))
				suite.Require().NoError(err)
			},
		},
		{
			"success: collateral USDC, trading asset ATOM, stop loss price 0, TakeProfitPrice 0",
			"",
			func() {
				err := suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, govtypes.ModuleName, positionCreator, sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, suite.GetAccountIssueAmount())))
				suite.Require().NoError(err)
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = amount
			},
		},
		{
			"success: collateral ATOM, trading asset ATOM, stop loss price 0, TakeProfitPrice 0",
			"",
			func() {
				tokensIn := sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, math.NewInt(1000_000_000)), sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000_000_000)))
				suite.AddLiquidity(ammPool, addr[3], tokensIn)
				msg.Creator = addr[2].String()
				msg.Collateral.Denom = ptypes.ATOM
				msg.Collateral.Amount = amount
				msg.Leverage = math.LegacyOneDec().MulInt64(2)
			},
		},
		{
			"collateral is USDC, trading asset is ATOM, amm pool has enough USDC but not enough ATOM",
			"negative pool amount after swap",
			func() {

				suite.ResetAndSetSuite(addr, true, amount.MulRaw(1000), math.NewInt(2))
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = amount
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			err := msg.ValidateBasic()
			suite.Require().NoError(err)
			_, err = suite.app.PerpetualKeeper.Open(suite.ctx, msg, false)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
