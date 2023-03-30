package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// AssetInfoKeyPrefix is the prefix to retrieve all AssetInfo
	AssetInfoKeyPrefix = "AssetInfo/value/"
)

// AssetInfoKey returns the store key to retrieve a AssetInfo from the index fields
func AssetInfoKey(
	denom string,
) []byte {
	var key []byte

	indexBytes := []byte(denom)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
