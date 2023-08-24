package types_test

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

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
				TokenOutDenom: "uusdc",
			},
		},
		TokenIn:           sdk.Coin{Denom: "uelys", Amount: sdk.NewInt(100)},
		TokenOutMinAmount: sdk.ZeroInt(),
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
				TokenInDenom: "uelys",
			},
			{
				PoolId:       2,
				TokenInDenom: "uusdc",
			},
		},
		TokenOut:         sdk.Coin{Denom: "uusdt", Amount: sdk.NewInt(100)},
		TokenInMaxAmount: sdk.ZeroInt(),
	}, 1)

	require.Contains(t, string(key), string(expectedKey))
}
