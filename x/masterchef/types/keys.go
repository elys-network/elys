package types

import (
	"encoding/binary"

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

	ParamsKey = "ParamsKey"

	PoolInfoKeyPrefix = "PoolInfo"

	ExternalIncentiveIndexKeyPrefix = "IndexExternalIncentive"

	ExternalIncentiveKeyPrefix = "ExternalIncentive"

	PoolRewardInfoKeyPrefix = "PoolRewardInfo"

	UserRewardInfoKeyPrefix = "UserRewardInfo"

	PoolRewardsAccumKeyPrefix = "PoolRewardsAccum"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func PoolInfoKey(poolId uint64) []byte {
	var key []byte

	poolIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(poolIdBytes, poolId)
	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

func ExternalIncentiveKey(incentiveId uint64) []byte {
	var key []byte

	incentiveIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(incentiveIdBytes, incentiveId)
	key = append(key, incentiveIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

func ExternalIncentiveIndex() []byte {
	var key []byte

	key = append(key, ExternalIncentiveIndexKeyPrefix...)
	return key
}

func PoolRewardInfoKey(poolId uint64, rewardDenom string) []byte {
	var key []byte

	poolIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(poolIdBytes, poolId)
	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)
	key = append(key, rewardDenom...)
	key = append(key, []byte("/")...)

	return key
}

func UserRewardInfoKey(user string, poolId uint64, rewardDenom string) []byte {
	var key []byte

	key = append(key, user...)
	key = append(key, []byte("/")...)
	poolIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(poolIdBytes, poolId)
	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)
	key = append(key, rewardDenom...)
	key = append(key, []byte("/")...)

	return key
}

func GetPoolRewardsAccumPrefix(poolId uint64) []byte {
	return append([]byte(PoolRewardsAccumKeyPrefix), sdk.Uint64ToBigEndian(uint64(poolId))...)
}

func GetPoolRewardsAccumKey(poolId uint64, timestamp uint64) []byte {
	return append(GetPoolRewardsAccumPrefix(poolId), sdk.Uint64ToBigEndian(timestamp)...)
}
