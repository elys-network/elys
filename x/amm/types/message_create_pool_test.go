package types_test

import (
	"errors"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestMsgCreatePool_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgCreatePool
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgCreatePool{
				Sender: "invalid_address",
				PoolParams: types.PoolParams{
					SwapFee:   sdkmath.LegacyZeroDec(),
					UseOracle: false,
					FeeDenom:  ptypes.BaseCurrency,
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "swap fee is negative, invalid params",
			msg: types.MsgCreatePool{
				Sender: sample.AccAddress(),
				PoolParams: types.PoolParams{
					SwapFee:   sdkmath.LegacyNewDec(-1),
					UseOracle: false,
					FeeDenom:  ptypes.BaseCurrency,
				},
				PoolAssets: []types.PoolAsset{
					{
						Token:                  sdk.NewCoin("uusdc", sdkmath.NewInt(10000000)),
						Weight:                 sdkmath.NewInt(10),
						ExternalLiquidityRatio: sdkmath.LegacyOneDec(),
					},
					{
						Token:                  sdk.NewCoin("uatom", sdkmath.NewInt(10000000)),
						Weight:                 sdkmath.NewInt(10),
						ExternalLiquidityRatio: sdkmath.LegacyOneDec(),
					},
				},
			},
			err: types.ErrNegativeSwapFee,
		},
		{
			name: "pool assets must be exactly two",
			msg: types.MsgCreatePool{
				Sender: sample.AccAddress(),
				PoolParams: types.PoolParams{
					SwapFee:   sdkmath.LegacyZeroDec(),
					UseOracle: false,
					FeeDenom:  ptypes.BaseCurrency,
				},
				PoolAssets: []types.PoolAsset{
					{
						Token:  sdk.NewCoin("uatom", sdkmath.NewInt(10000000)),
						Weight: sdkmath.NewInt(10),
					},
				},
			},
			err: types.ErrPoolAssetsMustBeTwo,
		},
		{
			name: "Invalid Pool Assets",
			msg: types.MsgCreatePool{
				Sender: sample.AccAddress(),
				PoolParams: types.PoolParams{
					SwapFee:   sdkmath.LegacyZeroDec(),
					UseOracle: false,
					FeeDenom:  ptypes.BaseCurrency,
				},
				PoolAssets: []types.PoolAsset{
					{
						Token:                  sdk.NewCoin("uusdc", sdkmath.NewInt(10000000)),
						Weight:                 sdkmath.NewInt(10),
						ExternalLiquidityRatio: sdkmath.LegacyOneDec(),
					},
					{
						Token:                  sdk.NewCoin("uatom", sdkmath.NewInt(10000000)),
						Weight:                 sdkmath.NewInt(-1),
						ExternalLiquidityRatio: sdkmath.LegacyOneDec(),
					},
				},
			},
			err: errors.New("invalid pool asset"),
		},
		{
			name: "valid address",
			msg: types.MsgCreatePool{
				Sender: sample.AccAddress(),
				PoolParams: types.PoolParams{
					SwapFee:   sdkmath.LegacyZeroDec(),
					UseOracle: false,
					FeeDenom:  ptypes.BaseCurrency,
				},
				PoolAssets: []types.PoolAsset{
					{
						Token:                  sdk.NewCoin("uusdc", sdkmath.NewInt(10000000)),
						Weight:                 sdkmath.NewInt(10),
						ExternalLiquidityRatio: sdkmath.LegacyOneDec(),
					},
					{
						Token:                  sdk.NewCoin("uatom", sdkmath.NewInt(10000000)),
						Weight:                 sdkmath.NewInt(10),
						ExternalLiquidityRatio: sdkmath.LegacyOneDec(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestInitialLiquidity(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name          string
		msg           types.MsgCreatePool
		expectedCoins sdk.Coins
		expectedPanic bool
	}{
		{
			"successful initial liquidity with sorted poolAssets",
			types.MsgCreatePool{
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000))},
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000))},
				},
			},
			sdk.Coins{
				sdk.NewCoin("tokenA", sdkmath.NewInt(1000)),
				sdk.NewCoin("tokenB", sdkmath.NewInt(2000)),
			},
			false,
		},
		{
			"empty pool assets",
			types.MsgCreatePool{
				PoolAssets: []types.PoolAsset{},
			},
			sdk.Coins{},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedPanic {
				require.Panics(t, func() {
					_ = tc.msg.InitialLiquidity()
				})
			} else {
				coins := tc.msg.InitialLiquidity()
				require.Equal(t, tc.expectedCoins, coins)
			}
		})
	}
}
