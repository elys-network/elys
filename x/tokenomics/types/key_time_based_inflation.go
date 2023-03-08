package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// TimeBasedInflationKeyPrefix is the prefix to retrieve all TimeBasedInflation
	TimeBasedInflationKeyPrefix = "TimeBasedInflation/value/"
)

// TimeBasedInflationKey returns the store key to retrieve a TimeBasedInflation from the index fields
func TimeBasedInflationKey(
	startBlockHeight uint64,
	endBlockHeight uint64,
) []byte {
	var key []byte

	startBlockHeightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(startBlockHeightBytes, startBlockHeight)
	key = append(key, startBlockHeightBytes...)
	key = append(key, []byte("/")...)

	endBlockHeightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(endBlockHeightBytes, endBlockHeight)
	key = append(key, endBlockHeightBytes...)
	key = append(key, []byte("/")...)

	return key
}
