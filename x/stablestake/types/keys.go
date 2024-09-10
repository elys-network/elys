package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "stablestake"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_stablestake"
)

var (
	DebtPrefixKey     = []byte{0x01}
	InterestPrefixKey = []byte{0x02}
	ParamKeyPrefix    = []byte{0x03}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetDebtKey(owner sdk.AccAddress) []byte {
	return append(DebtPrefixKey, address.MustLengthPrefix(owner)...)
}

func GetLegacyDebtKey(owner string) []byte {
	return append(DebtPrefixKey, []byte(owner)...)
}
