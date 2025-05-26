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
	CommitmentsKeyPrefix               = []byte{0x01}
	ParamsKey                          = []byte{0x02}
	AtomStakersKeyPrefix               = []byte{0x03}
	NFTHoldersKeyPrefix                = []byte{0x04}
	CadetsKeyPrefix                    = []byte{0x05}
	GovernorKeyPrefix                  = []byte{0x06}
	AirdropClaimedKeyPrefix            = []byte{0x07}
	TotalClaimedKeyPrefix              = []byte{0x08}
	KolKeyPrefix                       = []byte{0x09}
	TotalSupplyKeyPrefix               = []byte{0x10}
	MajorKeyPrefix                     = []byte{0x11}
	RewardProgramKeyPrefix             = []byte{0x12}
	RewardProgramClaimedKeyPrefix      = []byte{0x13}
	TotalRewardProgramClaimedKeyPrefix = []byte{0x14}
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

func GetAirdropClaimedKey(addr sdk.AccAddress) []byte {
	return append(AirdropClaimedKeyPrefix, address.MustLengthPrefix(addr)...)
}

func GetkolKey(addr sdk.AccAddress) []byte {
	return append(KolKeyPrefix, address.MustLengthPrefix(addr)...)
}

func GetMajorKey(addr sdk.AccAddress) []byte {
	return append(MajorKeyPrefix, address.MustLengthPrefix(addr)...)
}

func GetRewardProgramKey(addr sdk.AccAddress) []byte {
	return append(RewardProgramKeyPrefix, address.MustLengthPrefix(addr)...)
}

func GetRewardProgramClaimedKey(addr sdk.AccAddress) []byte {
	return append(RewardProgramClaimedKeyPrefix, address.MustLengthPrefix(addr)...)
}
