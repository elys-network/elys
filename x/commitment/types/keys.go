package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "commitment"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_commitment"

	LegacyParamsKey            = "Params/value/"
	LegacyCommitmentsKeyPrefix = "Commitments/value/"
)

var (
	CommitmentsKeyPrefix = []byte{0x01}
	ParamsKey            = []byte{0x02}
)

func LegacyKeyPrefix(p string) []byte {
	return []byte(p)
}

func LegacyCommitmentsKey(
	creator string,
) []byte {
	var key []byte

	creatorBytes := []byte(creator)
	key = append(key, creatorBytes...)
	key = append(key, []byte("/")...)

	return key
}

func GetCommitmentsKey(creator sdk.AccAddress) []byte {
	return append(CommitmentsKeyPrefix, address.MustLengthPrefix(creator)...)
}
