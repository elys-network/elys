package types

import (
	"encoding/binary"
	fmt "fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "amm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// TStoreKey defines the transient store key
	TStoreKey = "transient_amm"

	// PoolKeyPrefix is the prefix to retrieve all Pool
	PoolKeyPrefix = "Pool/value/"
	// OraclePoolSlippageTrackPrefix is the prefix to retrieve slippage tracked
	OraclePoolSlippageTrackPrefix = "OraclePool/slippage/track/value/"

	// DenomLiquidityKeyPrefix is the prefix to retrieve all DenomLiquidity
	DenomLiquidityKeyPrefix = "DenomLiquidity/value/"

	TLastSwapRequestIndex  = "last-swap-request-index"
	TSwapExactAmountInKey  = "batch/swap-exact-amount-in"
	TSwapExactAmountOutKey = "batch/swap-exact-amount-out"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// PoolKey returns the store key to retrieve a Pool from the index fields
func PoolKey(poolId uint64) []byte {
	var key []byte

	poolIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(poolIdBytes, poolId)
	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

func OraclePoolSlippageTrackKey(poolId uint64, timestamp uint64) []byte {
	return append(sdk.Uint64ToBigEndian(poolId), sdk.Uint64ToBigEndian(timestamp)...)
}

// DenomLiquidityKey returns the store key to retrieve a DenomLiquidity from the index fields
func DenomLiquidityKey(denom string) []byte {
	var key []byte

	denomBytes := []byte(denom)
	key = append(key, denomBytes...)
	key = append(key, []byte("/")...)

	return key
}

func TKeyPrefixSwapExactAmountInPrefix(m *MsgSwapExactAmountIn) []byte {
	prefix := []byte(fmt.Sprintf("%s/", m.TokenIn.Denom))
	routeKeys := []string{}
	for _, route := range m.Routes[:1] {
		routeKeys = append(routeKeys, fmt.Sprintf("%d/%s", route.PoolId, route.TokenOutDenom))
	}
	prefix = append(prefix, []byte(strings.Join(routeKeys, "/"))...)
	return prefix
}

func FirstPoolIdFromSwapExactAmountIn(m *MsgSwapExactAmountIn) uint64 {
	if len(m.Routes) > 0 {
		return m.Routes[0].PoolId
	}
	return 0
}

func TKeyPrefixSwapExactAmountIn(m *MsgSwapExactAmountIn, index uint64) []byte {
	prefix := TKeyPrefixSwapExactAmountInPrefix(m)
	return append(prefix, sdk.Uint64ToBigEndian(index)...)
}

func TKeyPrefixSwapExactAmountOutPrefix(m *MsgSwapExactAmountOut) []byte {
	prefix := []byte(fmt.Sprintf("/%s", m.TokenOut.Denom))
	routeKeys := []string{}
	if len(m.Routes) > 0 {
		route := m.Routes[len(m.Routes)-1]
		routeKeys = append(routeKeys, fmt.Sprintf("%s/%d", route.TokenInDenom, route.PoolId))
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
