package types

import (
	"encoding/binary"
	"fmt"
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

	WeightAndSlippagePrefix = "WeightBreakingFee/slippage/value/"

	// DenomLiquidityKeyPrefix is the prefix to retrieve all DenomLiquidity
	DenomLiquidityKeyPrefix = "DenomLiquidity/value/"

	TLastSwapRequestIndex  = "last-swap-request-index"
	TSwapExactAmountInKey  = "batch/swap-exact-amount-in"
	TSwapExactAmountOutKey = "batch/swap-exact-amount-out"

	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = "Params/value/"
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

func WeightAndSlippageFeeKey(poolId uint64, date string) []byte {
	dateBytes := []byte(date)
	return append(sdk.Uint64ToBigEndian(poolId), dateBytes...)
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
	prefix := fmt.Appendf(nil, "%s/", m.TokenIn.Denom)
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
	sender := sdk.MustAccAddressFromBech32(m.Sender)
	// We prefix with address bytes here so that SwapRequest becomes deterministically random to prevent MEV attacks.
	// We do not add address length prefix here, as 32 bytes address might suffer as then they always will be in the end
	return append(prefix, append(sender, sdk.Uint64ToBigEndian(index)...)...)
}

func TKeyPrefixSwapExactAmountOutPrefix(m *MsgSwapExactAmountOut) []byte {
	prefix := fmt.Appendf(nil, "/%s", m.TokenOut.Denom)
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
	sender := sdk.MustAccAddressFromBech32(m.Sender)
	// We prefix with address bytes here so that SwapRequest becomes deterministically random to prevent MEV attacks.
	// We do not add address length prefix here, as 32 bytes address might suffer as then they always will be in the end
	return append(prefix, append(sender, sdk.Uint64ToBigEndian(index)...)...)
}
