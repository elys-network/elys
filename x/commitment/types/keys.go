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
)

const MaxPageLimit = 10000

var (
	CommitmentsKeyPrefix = []byte{0x01}
	ParamsKey            = []byte{0x02}
	AtomStakersKeyPrefix = []byte{0x03}
	NFTHoldersKeyPrefix  = []byte{0x04}
	CadetsKeyPrefix      = []byte{0x05}
	GovernorKeyPrefix    = []byte{0x06}
)

func GetCommitmentsKey(creator sdk.AccAddress) []byte {
	return append(CommitmentsKeyPrefix, address.MustLengthPrefix(creator)...)
}

func GetAtomStakerKey(addr sdk.AccAddress) []byte {
	return append(AtomStakersKeyPrefix, address.MustLengthPrefix(addr)...)
}

func GetNFTHolderKey(addr sdk.AccAddress) []byte {
	return append(NFTHoldersKeyPrefix, address.MustLengthPrefix(addr)...)
}

func GetCadetKey(addr sdk.AccAddress) []byte {
	return append(CadetsKeyPrefix, address.MustLengthPrefix(addr)...)
}

func GetGovernorKey(addr sdk.AccAddress) []byte {
	return append(GovernorKeyPrefix, address.MustLengthPrefix(addr)...)
}
