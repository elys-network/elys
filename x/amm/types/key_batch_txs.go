package types

import (
	"encoding/binary"
	fmt "fmt"
	"strings"
)

var _ binary.ByteOrder

const (
	TLastSwapRequestIndex  = "last-swap-request-index"
	TSwapExactAmountInKey  = "batch/swap-exact-amount-in"
	TSwapExactAmountOutKey = "batch/swap-exact-amount-out"
)

func TKeyPrefixSwapExactAmountIn(m *MsgSwapExactAmountIn, index uint64) []byte {
	prefix := []byte(m.TokenIn.Denom + "/")
	routeKeys := []string{}
	for _, route := range m.Routes {
		routeKeys = append(routeKeys, fmt.Sprintf("%d/%s", route.PoolId, route.TokenOutDenom))
	}
	prefix = append(prefix, []byte(strings.Join(routeKeys, "/"))...)
	return prefix
}

func TKeyPrefixSwapExactAmountOut(m *MsgSwapExactAmountOut, index uint64) []byte {
	prefix := []byte(m.TokenOut.Denom + "/")
	routeKeys := []string{}
	for _, route := range m.Routes {
		routeKeys = append(routeKeys, fmt.Sprintf("%d/%s", route.PoolId, route.TokenInDenom))
	}
	prefix = append(prefix, []byte(strings.Join(routeKeys, "/"))...)
	return prefix
}
