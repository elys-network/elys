package types

import (
	"encoding/binary"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "tradeshield"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

var (
	ParamsKey                     = []byte{0x01}
	PendingSpotOrderKey           = []byte{0x02}
	PendingSpotOrderCountKey      = []byte{0x03}
	SortedSpotOrderKey            = []byte{0x04}
	PendingPerpetualOrderKey      = []byte{0x05}
	PendingPerpetualOrderCountKey = []byte{0x06}
	SortedPerpetualOrderKey       = []byte{0x07}
)

func DecodeUint64Slice(bz []byte) ([]uint64, error) {
	if len(bz)%8 != 0 {
		return nil, errors.New("invalid byte slice length")
	}
	slice := make([]uint64, len(bz)/8)
	for i := range slice {
		slice[i] = binary.BigEndian.Uint64(bz[i*8:])
	}
	return slice, nil
}

func GetPendingPerpetualOrderAddressKey(user sdk.AccAddress) []byte {
	key := PendingPerpetualOrderKey
	key = append(key, address.MustLengthPrefix(user)...)
	return key
}

// GetPendingPerpetualOrderKeyBytes returns the byte representation of the Address + PoolId + OrderId
func GetPendingPerpetualOrderKeyBytes(user sdk.AccAddress, poolId uint64, orderId uint64) []byte {

	key := address.MustLengthPrefix(user)
	key = append(key, []byte("/")...)
	poolIdBytes := sdk.Uint64ToBigEndian(poolId)
	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)
	orderIdBytes := sdk.Uint64ToBigEndian(orderId)
	key = append(key, orderIdBytes...)

	return key
}

func GetPendingPerpetualOrderAddressPoolKey(user sdk.AccAddress, poolId uint64) []byte {
	key := address.MustLengthPrefix(user)
	key = append(key, []byte("/")...)
	poolIdBytes := sdk.Uint64ToBigEndian(poolId)
	key = append(key, poolIdBytes...)

	return key
}
