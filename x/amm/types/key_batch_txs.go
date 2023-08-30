package types

import (
	"encoding/binary"
	fmt "fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ binary.ByteOrder

const (
	TLastSwapRequestIndex  = "last-swap-request-index"
	TSwapExactAmountInKey  = "batch/swap-exact-amount-in"
	TSwapExactAmountOutKey = "batch/swap-exact-amount-out"
)

func TKeyPrefixSwapExactAmountInPrefix(m *MsgSwapExactAmountIn) []byte {
	prefix := []byte(m.TokenIn.Denom + "/")
	routeKeys := []string{}
	for _, route := range m.Routes[:1] {
		routeKeys = append(routeKeys, fmt.Sprintf("%d/%s", route.PoolId, route.TokenOutDenom))
	}
	prefix = append(prefix, []byte(strings.Join(routeKeys, "/"))...)
	return prefix
}

func FirstPoolIdFromSwapExactAmountIn(m *MsgSwapExactAmountIn) uint64 {
	for _, route := range m.Routes {
		return route.PoolId
	}
	return 0
}

func TKeyPrefixSwapExactAmountIn(m *MsgSwapExactAmountIn, index uint64) []byte {
	prefix := TKeyPrefixSwapExactAmountInPrefix(m)
	return append(prefix, sdk.Uint64ToBigEndian(index)...)
}

func TKeyPrefixSwapExactAmountOutPrefix(m *MsgSwapExactAmountOut) []byte {
	prefix := []byte("/" + m.TokenOut.Denom)
	routeKeys := []string{}
	for i := len(m.Routes) - 1; i >= 0; i-- {
		route := m.Routes[i]
		routeKeys = append(routeKeys, fmt.Sprintf("%s/%d", route.TokenInDenom, route.PoolId))
		break
	}
	prefix = append([]byte(strings.Join(routeKeys, "/")), prefix...)
	return prefix
}

func FirstPoolIdFromSwapExactAmountOut(m *MsgSwapExactAmountOut) uint64 {
	for i := len(m.Routes) - 1; i >= 0; i-- {
		route := m.Routes[i]
		return route.PoolId
	}
	return 0
}

func TKeyPrefixSwapExactAmountOut(m *MsgSwapExactAmountOut, index uint64) []byte {
	prefix := TKeyPrefixSwapExactAmountOutPrefix(m)
	return append(prefix, sdk.Uint64ToBigEndian(index)...)
}
