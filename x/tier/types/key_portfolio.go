package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// PortfolioKeyPrefix is the prefix to retrieve all Portfolio
	PortfolioKeyPrefix = "Portfolio/value/"
)

// PortfolioKey returns the store key to retrieve a Portfolio from the index fields
func PortfolioKey(
	asset string,
) []byte {
	var key []byte

	assetBytes := []byte(asset)
	key = append(key, assetBytes...)
	key = append(key, []byte("/")...)

	return key
}
