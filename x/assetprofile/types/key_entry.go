package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// EntryKeyPrefix is the prefix to retrieve all Entry
	EntryKeyPrefix = "Entry/value/"
)

// EntryKey returns the store key to retrieve a Entry from the index fields
func EntryKey(
	baseDenom string,
) []byte {
	var key []byte

	baseDenomBytes := []byte(baseDenom)
	key = append(key, baseDenomBytes...)
	key = append(key, []byte("/")...)

	return key
}
