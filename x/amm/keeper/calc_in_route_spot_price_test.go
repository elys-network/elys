package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/elys-network/elys/v5/x/amm/types"
	oracletypes "github.com/elys-network/elys/v5/x/oracle/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *AmmKeeperTestSuite) TestCalcInRouteSpotPrice() {
	poolInitBalance := sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)}
	pool2InitBalance := sdk.Coins{sdk.NewInt64Coin("uusda", 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)}
	senderInitBalance := sdk.Coins{sdk.NewInt64Coin(ptypes.Elys, 1000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 1000000)}

	suite.SetupTest()
	suite.SetupStableCoinPrices()

	// Set up oracle asset info
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:   ptypes.Elys,
		Decimal: 6,
	})
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:   ptypes.BaseCurrency,
		Decimal: 6,
	})
	suite.app.OracleKeeper.SetAssetInfo(suite.ctx, oracletypes.AssetInfo{
		Denom:   "uusda",
		Decimal: 6,
	})

	// bootstrap accounts
	sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	poolAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	treasury2Addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	// bootstrap balances
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, senderInitBalance)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, senderInitBalance)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, poolInitBalance)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, poolInitBalance)
	suite.Require().NoError(err)

	// execute function
	for _, coin := range poolInitBalance {
		suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, types.DenomLiquidity{
			Denom:     coin.Denom,
			Liquidity: coin.Amount,
		})
	}
	pool := types.Pool{
		PoolId:            1,
		Address:           types.NewPoolAddress(uint64(1)).String(),
		RebalanceTreasury: treasuryAddr.String(),
		PoolParams: types.PoolParams{
			UseOracle: false,
			SwapFee:   sdkmath.LegacyZeroDec(),
			FeeDenom:  ptypes.BaseCurrency,
		},
		TotalShares: sdk.Coin{},
		PoolAssets: []types.PoolAsset{
			{
				Token:  poolInitBalance[0],
				Weight: sdkmath.NewInt(10),
			},
			{
				Token:  poolInitBalance[1],
				Weight: sdkmath.NewInt(10),
			},
		},
		TotalWeight: sdkmath.ZeroInt(),
	}
	pool2 := types.Pool{
		PoolId:            2,
		Address:           types.NewPoolAddress(uint64(2)).String(),
		RebalanceTreasury: treasury2Addr.String(),
		PoolParams: types.PoolParams{
			SwapFee:  sdkmath.LegacyZeroDec(),
			FeeDenom: ptypes.BaseCurrency,
		},
		TotalShares: sdk.Coin{},
		PoolAssets: []types.PoolAsset{
			{
				Token:  pool2InitBalance[0],
				Weight: sdkmath.NewInt(10),
			},
			{
				Token:  pool2InitBalance[1],
				Weight: sdkmath.NewInt(10),
			},
		},
		TotalWeight: sdkmath.ZeroInt(),
	}
	suite.app.AmmKeeper.SetPool(suite.ctx, pool)
	suite.app.AmmKeeper.SetPool(suite.ctx, pool2)

	tokenIn := sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(100))
	routes := []*types.SwapAmountInRoute{{PoolId: 1, TokenOutDenom: ptypes.BaseCurrency}}
	spotPrice, _, _, totalDiscountedSwapFee, _, _, _, _, err := suite.app.AmmKeeper.CalcInRouteSpotPrice(suite.ctx, tokenIn, routes, osmomath.ZeroBigDec(), osmomath.MustNewBigDecFromStr("0.1"))
	suite.Require().NoError(err)
	suite.Require().Equal(spotPrice.String(), osmomath.OneBigDec().String())
	suite.Require().Equal(sdkmath.LegacyMustNewDecFromStr("0.1"), totalDiscountedSwapFee.Dec())

	routes = []*types.SwapAmountInRoute{
		{PoolId: 1, TokenOutDenom: ptypes.BaseCurrency},
		{PoolId: 2, TokenOutDenom: "uusda"},
	}
	spotPrice, _, _, _, _, _, _, _, err = suite.app.AmmKeeper.CalcInRouteSpotPrice(suite.ctx, tokenIn, routes, osmomath.ZeroBigDec(), osmomath.ZeroBigDec())
	suite.Require().NoError(err)
	suite.Require().Equal(spotPrice.String(), osmomath.OneBigDec().String())

	// Test no routes
	_, _, _, _, _, _, _, _, err = suite.app.AmmKeeper.CalcInRouteSpotPrice(suite.ctx, tokenIn, nil, osmomath.ZeroBigDec(), osmomath.ZeroBigDec())
	suite.Require().Error(err)

	// Test invalid pool
	routes = []*types.SwapAmountInRoute{{PoolId: 9999, TokenOutDenom: "uusda"}}
	_, _, _, _, _, _, _, _, err = suite.app.AmmKeeper.CalcInRouteSpotPrice(suite.ctx, tokenIn, routes, osmomath.ZeroBigDec(), osmomath.ZeroBigDec())
	suite.Require().Error(err)
}
