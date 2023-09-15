package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// AccountedPoolKeyPrefix is the prefix to retrieve all AccountedPool
	AccountedPoolKeyPrefix = "AccountedPool/value/"
)

// AccountedPoolKey returns the store key to retrieve a AccountedPool from the index fields
func AccountedPoolKey(
	poolId uint64,
) []byte {
	var key []byte

	poolIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(poolIdBytes, poolId)
	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
