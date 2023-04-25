package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ElysDelegatorKeyPrefix is the prefix to retrieve all ElysDelegator
	ElysDelegatorKeyPrefix = "ElysDelegator/value/"
)

// ElysDelegatorKey returns the store key to retrieve a ElysDelegator from the index fields
func ElysDelegatorKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
