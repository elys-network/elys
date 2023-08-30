package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// DenomLiquidityKeyPrefix is the prefix to retrieve all DenomLiquidity
	DenomLiquidityKeyPrefix = "DenomLiquidity/value/"
)

// DenomLiquidityKey returns the store key to retrieve a DenomLiquidity from the index fields
func DenomLiquidityKey(denom string) []byte {
	var key []byte

	denomBytes := []byte(denom)
	key = append(key, denomBytes...)
	key = append(key, []byte("/")...)

	return key
}
