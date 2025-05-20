package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/v4/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/v4/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
	"github.com/elys-network/elys/v4/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *PerpetualKeeperTestSuite) TestOpenShort() {
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
	msg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(2),
		Position:        types.Position_SHORT,
		PoolId:          poolId,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: math.LegacyZeroDec(),
		StopLossPrice:   math.LegacyZeroDec(),
	}
	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func()
		postValidateFunction func(mtp *types.MTP)
	}{
		{
			"pool not found",
			types.ErrPoolDoesNotExist.Error(),

			func() {
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"amm pool not found",
			"pool does not exist",

			func() {
				ammPool = suite.CreateNewAmmPool(poolCreator, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
				poolId = ammPool.PoolId
				enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Pool: leveragelpmoduletypes.AddPool{
						poolId,
						math.LegacyMustNewDecFromStr("10"),
					},
				}
				_, err := leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
				suite.Require().NoError(err)

				suite.app.AmmKeeper.RemovePool(suite.ctx, ammPool.PoolId)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"trading asset is not in the pool",
			"(uelys) does not exist in the pool",

			func() {
				suite.app.AmmKeeper.SetPool(suite.ctx, ammPool)
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.TradingAsset = ptypes.Elys
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral asset neither base currency nor present in the pool",
			"collateral must be base currency",

			func() {
				msg.Collateral.Denom = ptypes.Elys
				msg.TradingAsset = ptypes.ATOM
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral asset is ATOM",
			"collateral must be base currency",

			func() {
				suite.app.AmmKeeper.SetPool(suite.ctx, ammPool)
				msg.Collateral.Denom = ptypes.ATOM
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"Borrow fails: lack of funds",
			"user does not have enough balance of the required coin",

			func() {
				msg.Collateral.Denom = ptypes.BaseCurrency
				params = suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.BorrowInterestRateMin = math.LegacyMustNewDecFromStr("0.12")
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, positionCreator, govtypes.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, suite.GetAccountIssueAmount())))
				suite.Require().NoError(err)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral amount is too high",
			"borrowed amount is higher than pool depth",

			func() {
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = msg.Collateral.Amount.MulRaw(1000_000_000)
			},
			func(mtp *types.MTP) {
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
				msg.TradingAsset = ptypes.ATOM
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral is USDC, trading asset is ATOM, amm pool has enough USDC but not enough ATOM",
			"amount too low",

			func() {
				suite.ResetAndSetSuite(addr, true, amount.MulRaw(1000), math.NewInt(2))
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = amount
				msg.TradingAsset = ptypes.ATOM
			},
			func(mtp *types.MTP) {
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			err := msg.ValidateBasic()
			suite.Require().NoError(err)
			mtp, err := suite.app.PerpetualKeeper.OpenDefineAssets(suite.ctx, poolId, msg, ptypes.BaseCurrency)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
			tc.postValidateFunction(mtp)
		})
	}
}
