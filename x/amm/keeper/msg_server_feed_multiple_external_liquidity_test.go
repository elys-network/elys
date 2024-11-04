package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestLiquidityRatioFromPriceDepth() {
	depth := sdk.NewDecWithPrec(1, 2) // 1%
	suite.Require().Equal("0.005012562893380045", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdk.NewDecWithPrec(2, 2) // 2%
	suite.Require().Equal("0.010050506338833466", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdk.NewDecWithPrec(5, 2) // 5%
	suite.Require().Equal("0.025320565519103609", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdk.NewDecWithPrec(10, 2) // 10%
	suite.Require().Equal("0.051316701949486200", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdk.NewDecWithPrec(30, 2) // 30%
	suite.Require().Equal("0.163339973465924452", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdk.NewDecWithPrec(50, 2) // 50%
	suite.Require().Equal("0.292893218813452475", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdk.NewDecWithPrec(70, 2) // 70%
	suite.Require().Equal("0.452277442494833886", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdk.NewDecWithPrec(90, 2) // 90%
	suite.Require().Equal("0.683772233983162067", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdk.NewDecWithPrec(100, 2) // 100%
	suite.Require().Equal("1.000000000000000000", keeper.LiquidityRatioFromPriceDepth(depth).String())
}

func (suite *KeeperTestSuite) TestGetExternalLiquidityRatio() {
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
						Token:                  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
						Weight:                 sdk.NewInt(50),
						ExternalLiquidityRatio: sdk.NewDec(1),
					},
					{
						Token:                  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100000000)),
						Weight:                 sdk.NewInt(50),
						ExternalLiquidityRatio: sdk.NewDec(1),
					},
				},
			},
			amountDepthInfo: []types.AssetAmountDepth{
				{
					Asset:  "USDC",
					Amount: sdk.NewDec(1000000000),
					Depth:  sdk.MustNewDecFromStr("0.5"),
				},
				{
					Asset:  "ATOM",
					Amount: sdk.NewDec(1000000000),
					Depth:  sdk.MustNewDecFromStr("0.5"),
				},
			},
			expectedResult: []types.PoolAsset{
				{
					Token:                  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
					Weight:                 sdk.NewInt(50),
					ExternalLiquidityRatio: sdk.MustNewDecFromStr("34.142135623730950558"),
				},
				{
					Token:                  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100000000)),
					Weight:                 sdk.NewInt(50),
					ExternalLiquidityRatio: sdk.MustNewDecFromStr("34.142135623730950558"),
				},
			},
			expectedError: nil,
		},
		{
			name: "missing asset entry",
			pool: types.Pool{
				PoolAssets: []types.PoolAsset{
					{
						Token: sdk.NewCoin("asset1", sdk.NewInt(1000)),
					},
				},
			},
			amountDepthInfo: []types.AssetAmountDepth{
				{
					Asset:  "asset1",
					Amount: sdk.NewDec(500),
					Depth:  sdk.MustNewDecFromStr("0.5"),
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
						Token: sdk.NewCoin("asset2denom", sdk.NewInt(1000)),
					},
				},
			},
			amountDepthInfo: []types.AssetAmountDepth{
				{
					Asset:  "asset2",
					Amount: sdk.NewDec(500),
					Depth:  sdk.MustNewDecFromStr("0.5"),
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
						Token:                  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(0)),
						Weight:                 sdk.NewInt(50),
						ExternalLiquidityRatio: sdk.NewDec(1),
					},
					{
						Token:                  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100000000)),
						Weight:                 sdk.NewInt(50),
						ExternalLiquidityRatio: sdk.NewDec(1),
					},
				},
			},
			amountDepthInfo: []types.AssetAmountDepth{
				{
					Asset:  "USDC",
					Amount: sdk.NewDec(1000000000),
					Depth:  sdk.MustNewDecFromStr("0.5"),
				},
				{
					Asset:  "ATOM",
					Amount: sdk.NewDec(1000000000),
					Depth:  sdk.MustNewDecFromStr("0.5"),
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
						Token:                  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
						Weight:                 sdk.NewInt(50),
						ExternalLiquidityRatio: sdk.NewDec(1),
					},
					{
						Token:                  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100000000)),
						Weight:                 sdk.NewInt(50),
						ExternalLiquidityRatio: sdk.NewDec(1),
					},
				},
			},
			amountDepthInfo: []types.AssetAmountDepth{
				{
					Asset:  "USDC",
					Amount: sdk.NewDec(1000000000),
					Depth:  sdk.MustNewDecFromStr("0"),
				},
				{
					Asset:  "ATOM",
					Amount: sdk.NewDec(1000000000),
					Depth:  sdk.MustNewDecFromStr("0.5"),
				},
			},
			expectedResult: []types.PoolAsset{
				{
					Token:                  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
					Weight:                 sdk.NewInt(50),
					ExternalLiquidityRatio: sdk.MustNewDecFromStr("34.142135623730950558"),
				},
				{
					Token:                  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100000000)),
					Weight:                 sdk.NewInt(50),
					ExternalLiquidityRatio: sdk.MustNewDecFromStr("34.142135623730950558"),
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
