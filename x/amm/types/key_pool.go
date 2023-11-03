package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// PoolKeyPrefix is the prefix to retrieve all Pool
	PoolKeyPrefix = "Pool/value/"
	// OraclePoolSlippageTrackPrefix is the prefix to retrieve slippage tracked
	OraclePoolSlippageTrackPrefix = "OraclePool/slippage/track/value/"
)

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
