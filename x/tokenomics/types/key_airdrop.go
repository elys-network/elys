package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// AirdropKeyPrefix is the prefix to retrieve all Airdrop
	AirdropKeyPrefix = "Airdrop/value/"
)

// AirdropKey returns the store key to retrieve a Airdrop from the index fields
func AirdropKey(
	intent string,
) []byte {
	var key []byte

	intentBytes := []byte(intent)
	key = append(key, intentBytes...)
	key = append(key, []byte("/")...)

	return key
}
