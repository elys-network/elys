package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/v6/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/v6/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *PerpetualKeeperTestSuite) TestOpenShort() {
	suite.ResetSuite()
	suite.SetupCoinPrices()
	addr := suite.AddAccounts(10, nil)
	amount := math.NewInt(1000)
	poolCreator := addr[0]
	positionCreator := addr[1]
	poolId := uint64(1)
	var ammPool ammtypes.Pool
	params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
	params.BorrowInterestRateMin = math.LegacyMustNewDecFromStr("0.12")
	_ = suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)

	ammPool = suite.CreateNewAmmPool(poolCreator, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
	poolId = ammPool.PoolId
	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err := leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
	suite.Require().NoError(err)

	msg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(2),
		Position:        types.Position_SHORT,
		PoolId:          poolId,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyZeroDec(),
		StopLossPrice:   math.LegacyZeroDec(),
	}
	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{
			"collateral amount is too high",
			"borrowed amount is higher than pool depth",

			func() {
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = msg.Collateral.Amount.MulRaw(1000_000_000)
			},
		},
		{
			"success: collateral USDC, trading asset ATOM, stop loss price 0, TakeProfitPrice 0",
			"",

			func() {
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = amount
			},
		},
		{
			"collateral is USDC, trading asset is ATOM, amm pool has enough USDC but not enough ATOM",
			"mtp health would be too low for safety factor",
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
			_, err = suite.app.PerpetualKeeper.Open(suite.ctx, msg)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
