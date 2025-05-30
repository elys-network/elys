package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name
	ModuleName = "vaults"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_vaults"
)

var (
	ParamKeyPrefix    = []byte{0x01}
	VaultKeyPrefix    = []byte{0x02}
	UserDataKeyPrefix = []byte{0x03}
)

func GetVaultKey(key uint64) []byte {
	return append(VaultKeyPrefix, sdk.Uint64ToBigEndian(key)...)
}

func GetUserDataKey(key string) []byte {
	return append(UserDataKeyPrefix, []byte(key)...)
}
