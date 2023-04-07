package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// CommitmentsKeyPrefix is the prefix to retrieve all Commitments
	CommitmentsKeyPrefix = "Commitments/value/"
)

// CommitmentsKey returns the store key to retrieve a Commitments from the index fields
func CommitmentsKey(
	creator string,
) []byte {
	var key []byte

	baseDenomBytes := []byte(creator)
	key = append(key, baseDenomBytes...)
	key = append(key, []byte("/")...)

	return key
}
