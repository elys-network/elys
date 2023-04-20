package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// HistoryKeyPrefix is the prefix to retrieve all History
	HistoryKeyPrefix = "History/value/"
)

// HistoryKey returns the store key to retrieve a History from the index fields
func HistoryKey(
	timestamp string,
	denom string,
) []byte {
	var key []byte

	timestampBytes := []byte(timestamp)
	key = append(key, timestampBytes...)
	key = append(key, []byte("/")...)

	denomBytes := []byte(denom)
	key = append(key, denomBytes...)
	key = append(key, []byte("/")...)

	return key
}
