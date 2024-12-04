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
)

var (
	ParamsKey                       = []byte{0x01}
	UserRewardInfoKeyPrefix         = []byte{0x02}
	PoolInfoKeyPrefix               = []byte{0x03}
	ExternalIncentiveIndexKeyPrefix = []byte{0x04}
	ExternalIncentiveKeyPrefix      = []byte{0x05}
	PoolRewardInfoKeyPrefix         = []byte{0x06}
	PoolRewardsAccumKeyPrefix       = []byte{0x07}
	FeeInfoKeyPrefix                = []byte{0x08}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetPoolInfoKey(poolId uint64) []byte {
	key := PoolInfoKeyPrefix
	poolIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(poolIdBytes, poolId)
	key = append(key, poolIdBytes...)

	return key
}

func GetExternalIncentiveKey(incentiveId uint64) []byte {
	key := ExternalIncentiveKeyPrefix

	incentiveIdBytes := sdk.Uint64ToBigEndian(incentiveId)
	key = append(key, incentiveIdBytes...)

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

// Timestamp will be a date in the format of YYYY-MM-DD
func GetFeeInfoKey(timestamp string) []byte {
	key := FeeInfoKeyPrefix
	key = append(key, []byte("/")...)
	return append(key, []byte(timestamp)...)
}
