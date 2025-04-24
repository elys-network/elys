package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestOpen() {
	suite.SetupCoinPrices()
	addr := suite.AddAccounts(10, nil)
	amount := math.NewInt(1000)
	poolCreator := addr[0]
	positionCreator := addr[1]
	poolId := uint64(1)
	tradingAssetPrice, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
	params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
	suite.Require().NoError(err)

	var ammPool ammtypes.Pool
	msg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_LONG,
		PoolId:          poolId,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: tradingAssetPrice.MulInt64(8),
		StopLossPrice:   math.LegacyZeroDec(),
	}
	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func()
		postValidateFunction func(mtp *types.MTP)
	}{
		{
			"base currency not found",
			"asset profile not found for denom",

			func() {
				suite.app.AssetprofileKeeper.RemoveEntry(suite.ctx, ptypes.BaseCurrency)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"borrow asset is usdc in long",
			"invalid operation: the borrowed asset cannot be the base currency: invalid borrowing asset",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()

				msg.TradingAsset = ptypes.BaseCurrency
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"invalid collateral",
			"collateral must either match the borrowed asset or be the base currency: invalid borrowing asset",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()

				msg.Collateral.Denom = ptypes.ATOM
				msg.TradingAsset = ptypes.Elys
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"short base currency",
			"cannot take a short position against the base currency: invalid borrowing asset",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()

				msg.Position = types.Position_SHORT
				msg.TradingAsset = ptypes.BaseCurrency
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"short same coin as collateral",
			"invalid operation: collateral asset cannot be identical to the borrowed asset for a short position: invalid collateral asset",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()

				msg.Position = types.Position_SHORT
				msg.TradingAsset = ptypes.ATOM
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"short with nonUSDC coin",
			"invalid collateral: the collateral asset for a short position must be the base currency: invalid collateral asset",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()

				msg.Position = types.Position_SHORT
				msg.Collateral.Denom = ptypes.Elys
				msg.TradingAsset = ptypes.ATOM
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"user not whitelisted",
			"unauthorised: address not on whitelist",

			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
				suite.AddAccounts(10, addr)
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.WhitelistingEnabled = true
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				msg.Position = types.Position_LONG
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.TradingAsset = ptypes.ATOM
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"amm pool not found",
			"pool does not exist",

			func() {
				for _, account := range addr {
					suite.app.PerpetualKeeper.WhitelistAddress(suite.ctx, account)
				}

				ammPool = suite.CreateNewAmmPool(poolCreator, true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
				poolId = ammPool.PoolId
				enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Pool: leveragelpmoduletypes.AddPool{
						AmmPoolId:   poolId,
						LeverageMax: math.LegacyMustNewDecFromStr("10"),
					},
				}
				_, err = leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
				suite.Require().NoError(err)
				suite.app.AmmKeeper.RemovePool(suite.ctx, ammPool.PoolId)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral asset neither base currency nor present in the pool",
			"collateral must either match the borrowed asset or be the base currency",

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

			func() {
				msg.Collateral.Denom = ptypes.ATOM
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.BorrowInterestRateMin = math.LegacyZeroDec()
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
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
			"Borrow fails: lack of funds",
			"user does not have enough balance of the required coin",

			func() {
				msg.Collateral.Amount = amount
				msg.Leverage = math.LegacyMustNewDecFromStr("1.2")
				tokensIn := sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, math.NewInt(1000_000_000)), sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000_000_000)))
				suite.AddLiquidity(ammPool, addr[3], tokensIn)
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.BorrowInterestRateMin = math.LegacyMustNewDecFromStr("0.12")
				err = suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, positionCreator, govtypes.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, suite.GetAccountIssueAmount())))
				suite.Require().NoError(err)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"stop loss price is greater than current asset price for long",
			"stop loss price cannot be greater than equal to tradingAssetPrice for long",

			func() {
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, govtypes.ModuleName, positionCreator, sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, suite.GetAccountIssueAmount())))
				suite.Require().NoError(err)
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = amount
				msg.TradingAsset = ptypes.ATOM

				msg.Position = types.Position_LONG
				msg.StopLossPrice = math.LegacyOneDec().MulInt64(7)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"stop loss price is less than current asset price for short",
			"stop loss price cannot be less than equal to tradingAssetPrice for short",

			func() {
				msg.Position = types.Position_SHORT
				msg.StopLossPrice = math.LegacyOneDec().MulInt64(3)
				msg.TakeProfitPrice = math.LegacyOneDec().MulInt64(2)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"success: collateral USDC, trading asset ATOM, stop loss price 0, TakeProfitPrice 0",
			"",

			func() {
				msg.Position = types.Position_LONG
				msg.StopLossPrice = math.LegacyZeroDec()
				msg.TakeProfitPrice = math.LegacyOneDec().MulInt64(8)
			},
			func(mtp *types.MTP) {
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
				msg.TradingAsset = ptypes.ATOM
				msg.Leverage = math.LegacyOneDec().MulInt64(2)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral is USDC, trading asset is ATOM, amm pool has enough USDC but not enough ATOM",
			"negative pool amount after swap",

			func() {
				suite.ResetAndSetSuite(addr, true, amount.MulRaw(1000), math.NewInt(2))

				msg.Creator = positionCreator.String()
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = amount
				msg.TradingAsset = ptypes.ATOM
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"take profit price below minimum ratio",
			fmt.Sprintf("take profit price should be between %s and %s times of current market price for long", params.MinimumLongTakeProfitPriceRatio.String(), params.MaximumLongTakeProfitPriceRatio.String()),
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
				msg.TakeProfitPrice = tradingAssetPrice.Mul(params.MinimumLongTakeProfitPriceRatio).Quo(math.LegacyNewDec(2))
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"take profit price above maximum ratio",
			fmt.Sprintf("take profit price should be between %s and %s times of current market price for long", params.MinimumLongTakeProfitPriceRatio.String(), params.MaximumLongTakeProfitPriceRatio.String()),
			func() {
				msg.TakeProfitPrice = tradingAssetPrice.Mul(params.MaximumLongTakeProfitPriceRatio).Mul(math.LegacyNewDec(2))
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
			_, err = suite.app.PerpetualKeeper.Open(suite.ctx, msg)
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
