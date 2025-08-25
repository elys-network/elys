package keeper_test

import (
	"cosmossdk.io/math"
	"fmt"
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

func (suite *PerpetualKeeperTestSuite) TestOpen() {
	suite.ResetSuite()
	suite.SetupCoinPrices()
	addr := suite.AddAccounts(10, nil)
	amount := math.NewInt(1000)
	poolCreator := addr[0]
	positionCreator := addr[1]
	poolId := uint64(1)
	tradingAssetPrice, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
	params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
	suite.Require().NoError(err)

	var creatorBalance sdk.Coins

	var initialPoolBankBalance sdk.Coins
	var initialAccountedPoolBalance sdk.Coins

	var finalPoolBankBalance sdk.Coins
	var finalAccountedPoolBalance sdk.Coins

	var ammPool ammtypes.Pool
	msg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(5),
		Position:        types.Position_LONG,
		PoolId:          poolId,
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
			"perpetual pool does not exist",
			"perpetual pool does not exist",
			func() {
				suite.SetupCoinPrices()

			},
			func(mtp *types.MTP) {
			},
		},
		{
			"invalid position",
			types.ErrInvalidPosition.Error(),
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
				_, err = leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
				suite.Require().NoError(err)
				msg.Position = types.Position_UNSPECIFIED

			},
			func(mtp *types.MTP) {
			},
		},
		{
			"invalid take profit price",
			"take profit price should be between",
			func() {
				msg.Position = types.Position_LONG
				msg.TakeProfitPrice = tradingAssetPrice.QuoInt64(2)

			},
			func(mtp *types.MTP) {
			},
		},
		{
			"user not whitelisted",
			"unauthorised: address not on whitelist",

			func() {
				suite.AddAccounts(10, addr)
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.WhitelistingEnabled = true
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				msg.Position = types.Position_LONG
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.TakeProfitPrice = tradingAssetPrice.MulInt64(8)
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
				suite.app.AmmKeeper.RemovePool(suite.ctx, ammPool.PoolId)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral is same as trading asset but pool doesn't have enough depth",
			"borrowed amount is higher than pool depth",

			func() {
				suite.app.AmmKeeper.SetPool(suite.ctx, ammPool)
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
				params = suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.BorrowInterestRateMin = math.LegacyMustNewDecFromStr("0.12")
				err = suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				creatorBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(msg.Creator))
				err = suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, sdk.MustAccAddressFromBech32(msg.Creator), govtypes.ModuleName, creatorBalance)
				suite.Require().NoError(err)
				enoughBalance := suite.app.BankKeeper.HasBalance(suite.ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Collateral)
				suite.Require().False(enoughBalance)

			},
			func(mtp *types.MTP) {
			},
		},
		{
			"success: collateral USDC, trading asset ATOM, stop loss price 0, TakeProfitPrice 0",
			"",

			func() {
				err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, govtypes.ModuleName, sdk.MustAccAddressFromBech32(msg.Creator), creatorBalance)
				suite.Require().NoError(err)
				enoughBalance := suite.app.BankKeeper.HasBalance(suite.ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Collateral)
				suite.Require().True(enoughBalance)

				msg.Position = types.Position_LONG
				msg.StopLossPrice = math.LegacyZeroDec()
				msg.TakeProfitPrice = math.LegacyOneDec().MulInt64(8)

				initialPoolBankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(ammPool.Address))
				accountedPool, found := suite.app.AccountedPoolKeeper.GetAccountedPool(suite.ctx, ammPool.PoolId)
				suite.Require().True(found)
				initialAccountedPoolBalance = accountedPool.TotalTokens
			},
			func(mtp *types.MTP) {
				finalPoolBankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(ammPool.Address))
				accountedPool, found := suite.app.AccountedPoolKeeper.GetAccountedPool(suite.ctx, ammPool.PoolId)
				suite.Require().True(found)
				finalAccountedPoolBalance = accountedPool.TotalTokens

				suite.Require().Equal(initialPoolBankBalance.Add(msg.Collateral), finalPoolBankBalance)
				suite.Require().Equal(initialAccountedPoolBalance.Add(msg.Collateral).Add(sdk.NewCoin(msg.Collateral.Denom, msg.Leverage.Sub(math.LegacyOneDec()).MulInt(msg.Collateral.Amount).TruncateInt())), finalAccountedPoolBalance)
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

				initialPoolBankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(ammPool.Address))
				accountedPool, found := suite.app.AccountedPoolKeeper.GetAccountedPool(suite.ctx, ammPool.PoolId)
				suite.Require().True(found)
				initialAccountedPoolBalance = accountedPool.TotalTokens
			},
			func(mtp *types.MTP) {
				finalPoolBankBalance = suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(ammPool.Address))
				accountedPool, found := suite.app.AccountedPoolKeeper.GetAccountedPool(suite.ctx, ammPool.PoolId)
				suite.Require().True(found)
				finalAccountedPoolBalance = accountedPool.TotalTokens

				suite.Require().Equal(initialPoolBankBalance.Add(msg.Collateral), finalPoolBankBalance)
				suite.Require().Equal(initialAccountedPoolBalance.Add(msg.Collateral).Add(sdk.NewCoin(msg.Collateral.Denom, msg.Leverage.Sub(math.LegacyOneDec()).MulInt(msg.Collateral.Amount).TruncateInt())), finalAccountedPoolBalance)
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
				pool := types.Pool{
					AmmPoolId:                            1,
					BaseAssetLiabilitiesRatio:            math.LegacyDec{},
					QuoteAssetLiabilitiesRatio:           math.LegacyDec{},
					BorrowInterestRate:                   math.LegacyDec{},
					PoolAssetsLong:                       []types.PoolAsset{{AssetDenom: ptypes.BaseCurrency}, {AssetDenom: ptypes.ATOM}},
					PoolAssetsShort:                      nil,
					LastHeightBorrowInterestRateComputed: 0,
					FundingRate:                          math.LegacyDec{},
					FeesCollected:                        nil,
					LeverageMax:                          math.LegacyDec{},
				}
				suite.app.PerpetualKeeper.SetPool(suite.ctx, pool)
				msg.TakeProfitPrice = tradingAssetPrice.Mul(params.MinimumLongTakeProfitPriceRatio).QuoInt64(2)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"take profit price above maximum ratio",
			fmt.Sprintf("take profit price should be between %s and %s times of current market price for long", params.MinimumLongTakeProfitPriceRatio.String(), params.MaximumLongTakeProfitPriceRatio.String()),
			func() {
				msg.TakeProfitPrice = tradingAssetPrice.Mul(params.MaximumLongTakeProfitPriceRatio).MulInt64(2)
			},
			func(mtp *types.MTP) {
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			_, err := suite.app.PerpetualKeeper.Open(suite.ctx, msg, false)
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
