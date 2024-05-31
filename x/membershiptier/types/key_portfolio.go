package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// PortfolioKeyPrefix is the prefix to retrieve all Portfolio
	PortfolioKeyPrefix = "Portfolio/value/"
)

// PortfolioKey returns the store key to retrieve a Portfolio from the index fields
func PortfolioKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
