package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "vaults"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_vaults"
)

var (
	ParamKeyPrefix            = []byte{0x01}
	VaultKeyPrefix            = []byte{0x02}
	UserDataKeyPrefix         = []byte{0x03}
	PoolInfoKeyPrefix         = []byte{0x04}
	PoolRewardInfoKeyPrefix   = []byte{0x05}
	UserRewardInfoKeyPrefix   = []byte{0x06}
	PoolRewardsAccumKeyPrefix = []byte{0x07}
	VaultPnLPrefix            = []byte{0x08}
)

func GetVaultKey(key uint64) []byte {
	return append(VaultKeyPrefix, sdk.Uint64ToBigEndian(key)...)
}

func GetUserDataKey(address string, vaultId uint64) []byte {
	key := append(UserDataKeyPrefix, []byte(address)...)
	key = append(key, []byte("/")...)
	key = append(key, sdk.Uint64ToBigEndian(vaultId)...)
	return key
}

func GetPoolInfoKey(key uint64) []byte {
	return append(PoolInfoKeyPrefix, sdk.Uint64ToBigEndian(key)...)
}

func GetPoolRewardInfoKey(poolId uint64, rewardDenom string) []byte {
	key := PoolRewardInfoKeyPrefix

	poolIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(poolIdBytes, poolId)
	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)
	key = append(key, rewardDenom...)

	return key
}

func GetUserRewardInfoKey(user sdk.AccAddress, poolId uint64, rewardDenom string) []byte {
	key := UserRewardInfoKeyPrefix

	key = append(key, address.MustLengthPrefix(user)...)
	key = append(key, []byte("/")...)
	poolIdBytes := sdk.Uint64ToBigEndian(poolId)
	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)
	key = append(key, rewardDenom...)

	return key
}

func GetPoolRewardsAccumPrefix(poolId uint64) []byte {
	key := PoolRewardsAccumKeyPrefix
	key = append(key, []byte("/")...)
	return append(key, sdk.Uint64ToBigEndian(poolId)...)
}

func GetPoolRewardsAccumKey(poolId uint64, timestamp uint64) []byte {
	key := GetPoolRewardsAccumPrefix(poolId)
	key = append(key, []byte("/")...)
	return append(key, sdk.Uint64ToBigEndian(timestamp)...)
}

func GetVaultPnLKey(vaultId string, date string) []byte {
	key := []byte(vaultId)
	key = append(key, []byte("/")...)
	return append(key, []byte(date)...)
}
