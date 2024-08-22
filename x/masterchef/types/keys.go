package types

import (
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/types/address"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "masterchef"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_masterchef"

	LegacyParamsKeyPrefix                 = "ParamsKey"
	LegacyUserRewardInfoKeyPrefix         = "UserRewardInfo"
	LegacyPoolInfoKeyPrefix               = "PoolInfo"
	LegacyExternalIncentiveIndexKeyPrefix = "IndexExternalIncentive"
	LegacyExternalIncentiveKeyPrefix      = "ExternalIncentive"
	LegacyPoolRewardInfoKeyPrefix         = "PoolRewardInfo"
	LegacyPoolRewardsAccumKeyPrefix       = "PoolRewardsAccum"
)

var (
	ParamsKey                       = []byte{0x01}
	UserRewardInfoKeyPrefix         = []byte{0x02}
	PoolInfoKeyPrefix               = []byte{0x03}
	ExternalIncentiveIndexKeyPrefix = []byte{0x04}
	ExternalIncentiveKeyPrefix      = []byte{0x05}
	PoolRewardInfoKeyPrefix         = []byte{0x06}
	PoolRewardsAccumKeyPrefix       = []byte{0x07}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func LegacyPoolInfoKey(poolId uint64) []byte {
	var key []byte

	poolIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(poolIdBytes, poolId)
	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

func GetPoolInfoKey(poolId uint64) []byte {
	key := PoolInfoKeyPrefix
	poolIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(poolIdBytes, poolId)
	key = append(key, poolIdBytes...)

	return key
}

func LegacyExternalIncentiveKey(incentiveId uint64) []byte {
	var key []byte

	incentiveIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(incentiveIdBytes, incentiveId)
	key = append(key, incentiveIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

func GetExternalIncentiveKey(incentiveId uint64) []byte {
	key := ExternalIncentiveKeyPrefix

	incentiveIdBytes := sdk.Uint64ToBigEndian(incentiveId)
	key = append(key, incentiveIdBytes...)

	return key
}

func LegacyExternalIncentiveIndex() []byte {
	var key []byte

	key = append(key, LegacyExternalIncentiveIndexKeyPrefix...)
	return key
}

func LegacyPoolRewardInfoKey(poolId uint64, rewardDenom string) []byte {
	var key []byte

	poolIdBytes := sdk.Uint64ToBigEndian(poolId)
	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)
	key = append(key, rewardDenom...)
	key = append(key, []byte("/")...)

	return key
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

func GetLegacyUserRewardInfoKey(user string, poolId uint64, rewardDenom string) []byte {
	key := KeyPrefix(LegacyUserRewardInfoKeyPrefix)

	key = append(key, user...)
	key = append(key, []byte("/")...)
	poolIdBytes := sdk.Uint64ToBigEndian(poolId)
	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)
	key = append(key, rewardDenom...)
	key = append(key, []byte("/")...)

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

func GetLegacyPoolRewardsAccumPrefix(poolId uint64) []byte {
	return append([]byte(LegacyPoolRewardsAccumKeyPrefix), sdk.Uint64ToBigEndian(uint64(poolId))...)
}

func GetLegacyPoolRewardsAccumKey(poolId uint64, timestamp uint64) []byte {
	return append(GetLegacyPoolRewardsAccumPrefix(poolId), sdk.Uint64ToBigEndian(timestamp)...)
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
