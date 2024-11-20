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
)

var (
	DebtPrefixKey     = []byte{0x01}
	InterestPrefixKey = []byte{0x02}
	ParamKeyPrefix    = []byte{0x03}
)

func GetDebtKey(owner sdk.AccAddress) []byte {
	return append(DebtPrefixKey, address.MustLengthPrefix(owner)...)
}
