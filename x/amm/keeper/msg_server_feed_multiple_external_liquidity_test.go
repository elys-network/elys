package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func (suite *AmmKeeperTestSuite) TestLiquidityRatioFromPriceDepth() {
	depth := sdkmath.LegacyNewDecWithPrec(1, 2) // 1%
	suite.Require().Equal("0.005012562893380045", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(2, 2) // 2%
	suite.Require().Equal("0.010050506338833466", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(5, 2) // 5%
	suite.Require().Equal("0.025320565519103609", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(10, 2) // 10%
	suite.Require().Equal("0.051316701949486200", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(30, 2) // 30%
	suite.Require().Equal("0.163339973465924452", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(50, 2) // 50%
	suite.Require().Equal("0.292893218813452475", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(70, 2) // 70%
	suite.Require().Equal("0.452277442494833886", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(90, 2) // 90%
	suite.Require().Equal("0.683772233983162067", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(100, 2) // 100%
	suite.Require().Equal("1.000000000000000000", keeper.LiquidityRatioFromPriceDepth(depth).String())
}

func (suite *AmmKeeperTestSuite) TestGetExternalLiquidityRatio() {
	suite.SetupTest()
	suite.SetupCoinPrices()
	// set asset profile
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.ATOM,
		Denom:       ptypes.ATOM,
		Decimals:    6,
		DisplayName: "ATOM",
	})
	// set asset profile
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.BaseCurrency,
		Denom:       ptypes.BaseCurrency,
		Decimals:    6,
		DisplayName: "USDC",
	})
	// set asset profile
	suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
		BaseDenom:   "asset2denom",
		Denom:       "asset2denom",
		Decimals:    6,
		DisplayName: "asset2",
	})
	tests := []struct {
		name            string
		pool            types.Pool
		amountDepthInfo []types.AssetAmountDepth
		expectedResult  []types.PoolAsset
		expectedError   error
	}{
		{
			name: "valid inputs",
			pool: types.Pool{
				PoolAssets: []types.PoolAsset{
					{
						Token:                  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000000)),
						Weight:                 sdkmath.NewInt(50),
						ExternalLiquidityRatio: sdkmath.LegacyNewDec(1),
					},
					{
						Token:                  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100000000)),
						Weight:                 sdkmath.NewInt(50),
						ExternalLiquidityRatio: sdkmath.LegacyNewDec(1),
					},
				},
			},
			amountDepthInfo: []types.AssetAmountDepth{
				{
					Asset:  "USDC",
					Amount: sdkmath.LegacyNewDec(1000000000),
					Depth:  sdkmath.LegacyMustNewDecFromStr("0.5"),
				},
				{
					Asset:  "ATOM",
					Amount: sdkmath.LegacyNewDec(1000000000),
					Depth:  sdkmath.LegacyMustNewDecFromStr("0.5"),
				},
			},
			expectedResult: []types.PoolAsset{
				{
					Token:                  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000000)),
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyMustNewDecFromStr("34.142135623730950558"),
				},
				{
					Token:                  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100000000)),
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyMustNewDecFromStr("34.142135623730950558"),
				},
			},
			expectedError: nil,
		},
		{
			name: "missing asset entry",
			pool: types.Pool{
				PoolAssets: []types.PoolAsset{
					{
						Token: sdk.NewCoin("asset1", sdkmath.NewInt(1000)),
					},
				},
			},
			amountDepthInfo: []types.AssetAmountDepth{
				{
					Asset:  "asset1",
					Amount: sdkmath.LegacyNewDec(500),
					Depth:  sdkmath.LegacyMustNewDecFromStr("0.5"),
				},
			},
			expectedResult: nil,
			expectedError:  fmt.Errorf("asset profile not found for denom"),
		},
		{
			name: "missing asset price",
			pool: types.Pool{
				PoolAssets: []types.PoolAsset{
					{
						Token: sdk.NewCoin("asset2denom", sdkmath.NewInt(1000)),
					},
				},
			},
			amountDepthInfo: []types.AssetAmountDepth{
				{
					Asset:  "asset2",
					Amount: sdkmath.LegacyNewDec(500),
					Depth:  sdkmath.LegacyMustNewDecFromStr("0.5"),
				},
			},
			expectedResult: nil,
			expectedError:  fmt.Errorf("asset price not set: asset2"),
		},
		{
			name: "division by zero",
			pool: types.Pool{
				PoolAssets: []types.PoolAsset{
					{
						Token:                  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(0)),
						Weight:                 sdkmath.NewInt(50),
						ExternalLiquidityRatio: sdkmath.LegacyNewDec(1),
					},
					{
						Token:                  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100000000)),
						Weight:                 sdkmath.NewInt(50),
						ExternalLiquidityRatio: sdkmath.LegacyNewDec(1),
					},
				},
			},
			amountDepthInfo: []types.AssetAmountDepth{
				{
					Asset:  "USDC",
					Amount: sdkmath.LegacyNewDec(1000000000),
					Depth:  sdkmath.LegacyMustNewDecFromStr("0.5"),
				},
				{
					Asset:  "ATOM",
					Amount: sdkmath.LegacyNewDec(1000000000),
					Depth:  sdkmath.LegacyMustNewDecFromStr("0.5"),
				},
			},
			expectedResult: nil,
			expectedError:  types.ErrAmountTooLow,
		},
		{
			name: "liquidity ratio is zero",
			pool: types.Pool{
				PoolAssets: []types.PoolAsset{
					{
						Token:                  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000000)),
						Weight:                 sdkmath.NewInt(50),
						ExternalLiquidityRatio: sdkmath.LegacyNewDec(1),
					},
					{
						Token:                  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100000000)),
						Weight:                 sdkmath.NewInt(50),
						ExternalLiquidityRatio: sdkmath.LegacyNewDec(1),
					},
				},
			},
			amountDepthInfo: []types.AssetAmountDepth{
				{
					Asset:  "USDC",
					Amount: sdkmath.LegacyNewDec(1000000000),
					Depth:  sdkmath.LegacyMustNewDecFromStr("0"),
				},
				{
					Asset:  "ATOM",
					Amount: sdkmath.LegacyNewDec(1000000000),
					Depth:  sdkmath.LegacyMustNewDecFromStr("0.5"),
				},
			},
			expectedResult: []types.PoolAsset{
				{
					Token:                  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000000)),
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyMustNewDecFromStr("34.142135623730950558"),
				},
				{
					Token:                  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100000000)),
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyMustNewDecFromStr("34.142135623730950558"),
				},
			},
			expectedError: types.ErrAmountTooLow,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			assets, err := suite.app.AmmKeeper.GetExternalLiquidityRatio(suite.ctx, tt.pool, tt.amountDepthInfo)
			if tt.expectedError != nil {
				require.Error(suite.T(), err)
				require.EqualError(suite.T(), err, tt.expectedError.Error())
			} else {
				require.NoError(suite.T(), err)
				require.Equal(suite.T(), tt.expectedResult, assets)
			}
		})
	}
}

func (suite *AmmKeeperTestSuite) TestFeedMultipleExternalLiquidity() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"feed multiple exeternal liquidity with empty msg",
			func() {
				suite.ResetSuite()
			},
			func() {
				// msg server
				msgServer := keeper.NewMsgServerImpl(*suite.app.AmmKeeper)
				goCtx := sdk.WrapSDKContext(suite.ctx)

				_, err := msgServer.FeedMultipleExternalLiquidity(goCtx, &types.MsgFeedMultipleExternalLiquidity{})
				suite.Require().Error(err)
			},
		},
		{
			"feed multiple exeternal liquidity without price feeder",
			func() {
				suite.ResetSuite()
			},
			func() {
				// msg server
				msgServer := keeper.NewMsgServerImpl(*suite.app.AmmKeeper)
				goCtx := sdk.WrapSDKContext(suite.ctx)

				addr := suite.AddAccounts(1, nil)[0]

				_, err := msgServer.FeedMultipleExternalLiquidity(goCtx, &types.MsgFeedMultipleExternalLiquidity{
					Sender: addr.String(),
					Liquidity: []types.ExternalLiquidity{
						{
							PoolId: 1,
							AmountDepthInfo: []types.AssetAmountDepth{
								{
									Asset:  ptypes.BaseCurrency,
									Amount: sdkmath.LegacyNewDec(1000000000),
									Depth:  sdkmath.LegacyMustNewDecFromStr("0.5"),
								},
							},
						},
					},
				})
				suite.Require().Error(err)
			},
		},
		{
			"feed multiple exeternal liquidity with price feeder but not active",
			func() {
				suite.ResetSuite()
			},
			func() {
				// msg server
				msgServer := keeper.NewMsgServerImpl(*suite.app.AmmKeeper)
				goCtx := sdk.WrapSDKContext(suite.ctx)

				addr := suite.AddAccounts(1, nil)[0]

				suite.app.OracleKeeper.SetPriceFeeder(suite.ctx, oracletypes.PriceFeeder{
					Feeder:   addr.String(),
					IsActive: false,
				})

				_, err := msgServer.FeedMultipleExternalLiquidity(goCtx, &types.MsgFeedMultipleExternalLiquidity{
					Sender: addr.String(),
					Liquidity: []types.ExternalLiquidity{
						{
							PoolId: 1,
							AmountDepthInfo: []types.AssetAmountDepth{
								{
									Asset:  ptypes.BaseCurrency,
									Amount: sdkmath.LegacyNewDec(1000000000),
									Depth:  sdkmath.LegacyMustNewDecFromStr("0.5"),
								},
							},
						},
					},
				})
				suite.Require().Error(err)
			},
		},
		{
			"feed multiple exeternal liquidity with active price feeder but invalid pool id",
			func() {
				suite.ResetSuite()
			},
			func() {
				// msg server
				msgServer := keeper.NewMsgServerImpl(*suite.app.AmmKeeper)
				goCtx := sdk.WrapSDKContext(suite.ctx)

				addr := suite.AddAccounts(1, nil)[0]

				suite.app.OracleKeeper.SetPriceFeeder(suite.ctx, oracletypes.PriceFeeder{
					Feeder:   addr.String(),
					IsActive: true,
				})

				_, err := msgServer.FeedMultipleExternalLiquidity(goCtx, &types.MsgFeedMultipleExternalLiquidity{
					Sender: addr.String(),
					Liquidity: []types.ExternalLiquidity{
						{
							PoolId: 1,
							AmountDepthInfo: []types.AssetAmountDepth{
								{
									Asset:  ptypes.BaseCurrency,
									Amount: sdkmath.LegacyNewDec(1000000000),
									Depth:  sdkmath.LegacyMustNewDecFromStr("0.5"),
								},
							},
						},
					},
				})
				suite.Require().Error(err)
			},
		},
		{
			"feed multiple exeternal liquidity with active price feeder with pool",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				// msg server
				msgServer := keeper.NewMsgServerImpl(*suite.app.AmmKeeper)
				goCtx := sdk.WrapSDKContext(suite.ctx)

				addr := suite.AddAccounts(1, nil)[0]

				amount := math.NewInt(100000000000)
				pool := suite.CreateNewAmmPool(addr, true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))

				suite.app.OracleKeeper.SetPriceFeeder(suite.ctx, oracletypes.PriceFeeder{
					Feeder:   addr.String(),
					IsActive: true,
				})

				_, err := msgServer.FeedMultipleExternalLiquidity(goCtx, &types.MsgFeedMultipleExternalLiquidity{
					Sender: addr.String(),
					Liquidity: []types.ExternalLiquidity{
						{
							PoolId: pool.PoolId,
							AmountDepthInfo: []types.AssetAmountDepth{
								{
									Asset:  ptypes.ATOM,
									Amount: sdkmath.LegacyNewDec(1000000000),
									Depth:  sdkmath.LegacyMustNewDecFromStr("0.5"),
								},
							},
						},
					},
				})
				suite.Require().Error(err)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			tc.postValidateFunction()
		})
	}
}
