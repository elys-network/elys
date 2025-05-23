package keeper_test

import (
	"errors"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/osmosis-labs/osmosis/osmomath"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/amm/keeper"
	"github.com/elys-network/elys/v5/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	oracletypes "github.com/elys-network/elys/v5/x/oracle/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func (suite *AmmKeeperTestSuite) TestLiquidityRatioFromPriceDepth() {
	depth := osmomath.NewBigDecWithPrec(1, 2) // 1%
	suite.Require().Equal("0.005012562893380045265520178998793995", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = osmomath.NewBigDecWithPrec(2, 2) // 2%
	suite.Require().Equal("0.010050506338833465838817893053211345", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = osmomath.NewBigDecWithPrec(5, 2) // 5%
	suite.Require().Equal("0.025320565519103609316158680010039970", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = osmomath.NewBigDecWithPrec(10, 2) // 10%
	suite.Require().Equal("0.051316701949486200400331936670184440", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = osmomath.NewBigDecWithPrec(30, 2) // 30%
	suite.Require().Equal("0.163339973465924452021827974214812510", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = osmomath.NewBigDecWithPrec(50, 2) // 50%
	suite.Require().Equal("0.292893218813452475599155637895150960", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = osmomath.NewBigDecWithPrec(70, 2) // 70%
	suite.Require().Equal("0.452277442494833886543030217199197866", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = osmomath.NewBigDecWithPrec(90, 2) // 90%
	suite.Require().Equal("0.683772233983162066800110645556728146", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = osmomath.NewBigDecWithPrec(100, 2) // 100%
	suite.Require().Equal("1.000000000000000000000000000000000000", keeper.LiquidityRatioFromPriceDepth(depth).String())
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
					ExternalLiquidityRatio: sdkmath.LegacyMustNewDecFromStr("34.142135623730950488"),
				},
				{
					Token:                  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100000000)),
					Weight:                 sdkmath.NewInt(50),
					ExternalLiquidityRatio: sdkmath.LegacyMustNewDecFromStr("34.142135623730950488"),
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
			expectedError:  errors.New("asset profile not found for denom"),
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

				_, err := msgServer.FeedMultipleExternalLiquidity(suite.ctx, &types.MsgFeedMultipleExternalLiquidity{Sender: authtypes.NewModuleAddress("test").String()})
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

				addr := suite.AddAccounts(1, nil)[0]

				_, err := msgServer.FeedMultipleExternalLiquidity(suite.ctx, &types.MsgFeedMultipleExternalLiquidity{
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

				addr := suite.AddAccounts(1, nil)[0]

				suite.app.OracleKeeper.SetPriceFeeder(suite.ctx, oracletypes.PriceFeeder{
					Feeder:   addr.String(),
					IsActive: false,
				})

				_, err := msgServer.FeedMultipleExternalLiquidity(suite.ctx, &types.MsgFeedMultipleExternalLiquidity{
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

				addr := suite.AddAccounts(1, nil)[0]

				suite.app.OracleKeeper.SetPriceFeeder(suite.ctx, oracletypes.PriceFeeder{
					Feeder:   addr.String(),
					IsActive: true,
				})

				_, err := msgServer.FeedMultipleExternalLiquidity(suite.ctx, &types.MsgFeedMultipleExternalLiquidity{
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

				addr := suite.AddAccounts(1, nil)[0]

				amount := sdkmath.NewInt(100000000000)
				pool := suite.CreateNewAmmPool(addr, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))

				suite.app.OracleKeeper.SetPriceFeeder(suite.ctx, oracletypes.PriceFeeder{
					Feeder:   addr.String(),
					IsActive: true,
				})

				_, err := msgServer.FeedMultipleExternalLiquidity(suite.ctx, &types.MsgFeedMultipleExternalLiquidity{
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
