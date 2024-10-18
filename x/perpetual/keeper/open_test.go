package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestOpen() {
	suite.SetupCoinPrices()
	addr := suite.AddAccounts(10, nil)
	amount := sdk.NewInt(1000)
	poolCreator := addr[0]
	positionCreator := addr[1]
	poolId := uint64(1)
	tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPriceByDenom(suite.ctx, ptypes.ATOM)
	suite.Require().NoError(err)
	var ammPool ammtypes.Pool
	msg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_LONG,
		PoolId:          poolId,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: tradingAssetPrice.MulInt64(4),
		StopLossPrice:   sdk.ZeroDec(),
	}
	testCases := []struct {
		name                 string
		expectErrMsg         string
		isBroker             bool
		prerequisiteFunction func()
		postValidateFunction func(mtp *types.MTP)
	}{
		{
			"base currency not found",
			"asset profile not found for denom",
			false,
			func() {
				suite.app.AssetprofileKeeper.RemoveEntry(suite.ctx, ptypes.BaseCurrency)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"user not whitelisted",
			"unauthorised: address not on whitelist",
			false,
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
				suite.AddAccounts(10, addr)
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.WhitelistingEnabled = true
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"pool not found",
			types.ErrPoolDoesNotExist.Error(),
			false,
			func() {
				for _, account := range addr {
					suite.app.PerpetualKeeper.WhitelistAddress(suite.ctx, account)
				}
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"pool is disabled",
			"perpetual not enabled for pool",
			false,
			func() {
				ammPool = suite.SetAndGetAmmPool(poolCreator, poolId, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))

				msgServer := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
				leverageLpPool := leveragelpmoduletypes.NewPool(poolId)
				leverageLpPool.Enabled = true
				leverageLpPool.Closed = false
				suite.app.LeveragelpKeeper.SetPool(suite.ctx, leverageLpPool)
				enablePoolMsg := types.MsgEnablePool{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					PoolId:    1,
				}
				_, err := msgServer.EnablePool(suite.ctx, &enablePoolMsg)
				suite.Require().NoError(err)
				pool, _ := suite.app.PerpetualKeeper.GetPool(suite.ctx, poolId)
				pool.Enabled = false
				suite.app.PerpetualKeeper.SetPool(suite.ctx, pool)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"amm pool not found",
			"pool does not exist",
			false,
			func() {
				pool, found := suite.app.PerpetualKeeper.GetPool(suite.ctx, 1)
				suite.Require().True(found)
				pool.Enabled = true
				suite.app.PerpetualKeeper.SetPool(suite.ctx, pool)
				suite.app.AmmKeeper.RemovePool(suite.ctx, ammPool.PoolId)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral asset neither base currency nor present in the pool",
			"collateral must either match the borrowed asset or be the base currency",
			false,
			func() {
				suite.app.AmmKeeper.SetPool(suite.ctx, ammPool)
				msg.Collateral.Denom = ptypes.Elys
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral is same as trading asset but pool doesn't have enough depth",
			"borrowed amount is higher than pool depth",
			false,
			func() {
				msg.Collateral.Denom = ptypes.ATOM
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.BorrowInterestRateMin = sdk.ZeroDec()
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral amount is too high",
			"borrowed amount is higher than pool depth",
			false,
			func() {
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = msg.Collateral.Amount.MulRaw(1000_000_000)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"Borrow fails: lack of funds",
			"user does not have enough balance of the required coin",
			false,
			func() {
				msg.Collateral.Amount = amount
				msg.Leverage = math.LegacyMustNewDecFromStr("1.2")
				tokensIn := sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1000_000_000)), sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1000_000_000)))
				suite.AddLiquidity(ammPool, addr[3], tokensIn)
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.BorrowInterestRateMin = sdk.MustNewDecFromStr("0.12")
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, positionCreator, govtypes.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, suite.GetAccountIssueAmount())))
				suite.Require().NoError(err)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"success: collateral USDC, trading asset ATOM, stop loss price 0, TakeProfitPrice 0",
			"",
			false,
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
			"success: collateral ATOM, trading asset ATOM, stop loss price 0, TakeProfitPrice 0",
			"",
			false,
			func() {
				tokensIn := sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1000_000_000)), sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1000_000_000)))
				suite.AddLiquidity(ammPool, addr[3], tokensIn)
				msg.Creator = addr[2].String()
				msg.Collateral.Denom = ptypes.ATOM
				msg.Collateral.Amount = amount
				msg.TradingAsset = ptypes.ATOM
				msg.Leverage = sdk.OneDec().MulInt64(2)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral is USDC, trading asset is ATOM, amm pool has enough USDC but not enough ATOM",
			"amount too low",
			false,
			func() {
				suite.ResetAndSetSuite(addr, 1, true, amount.MulRaw(1000), sdk.NewInt(2))

				msg.Creator = positionCreator.String()
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
			_, err = suite.app.PerpetualKeeper.Open(suite.ctx, msg, tc.isBroker)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
			//tc.postValidateFunction(mtp)
		})
	}
}
