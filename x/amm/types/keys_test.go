package types_test

import (
	"encoding/binary"
	"testing"

	sdkmath "cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/amm/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestPoolKey(t *testing.T) {
	poolID := uint64(1234567890)

	expectedKey := make([]byte, 8)
	binary.BigEndian.PutUint64(expectedKey, poolID)
	expectedKey = append(expectedKey, []byte("/")...)

	resultKey := types.PoolKey(poolID)

	require.Equal(t, expectedKey, resultKey)
}

func TestDenomLiquidityKey(t *testing.T) {
	denom := "liquidityToken"

	expectedKey := []byte("liquidityToken/")

	resultKey := types.DenomLiquidityKey(denom)

	require.Equal(t, expectedKey, resultKey)
}

func TestTKeyPrefixSwapExactAmountIn(t *testing.T) {
	expectedKey := []byte("uelys/1/uusdt")

	sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	key := types.TKeyPrefixSwapExactAmountIn(&types.MsgSwapExactAmountIn{
		Sender: sender.String(),
		Routes: []types.SwapAmountInRoute{
			{
				PoolId:        1,
				TokenOutDenom: "uusdt",
			},
			{
				PoolId:        2,
				TokenOutDenom: ptypes.BaseCurrency,
			},
		},
		TokenIn:           sdk.Coin{Denom: ptypes.Elys, Amount: sdkmath.NewInt(100)},
		TokenOutMinAmount: sdkmath.ZeroInt(),
	}, 1)

	require.Contains(t, string(key), string(expectedKey))
}

func TestTKeyPrefixSwapExactAmountOut(t *testing.T) {
	expectedKey := []byte("uusdc/2/uusdt")
	sender := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	key := types.TKeyPrefixSwapExactAmountOut(&types.MsgSwapExactAmountOut{
		Sender: sender.String(),
		Routes: []types.SwapAmountOutRoute{
			{
				PoolId:       1,
				TokenInDenom: ptypes.Elys,
			},
			{
				PoolId:       2,
				TokenInDenom: ptypes.BaseCurrency,
			},
		},
		TokenOut:         sdk.Coin{Denom: "uusdt", Amount: sdkmath.NewInt(100)},
		TokenInMaxAmount: sdkmath.ZeroInt(),
	}, 1)

	require.Contains(t, string(key), string(expectedKey))
}
