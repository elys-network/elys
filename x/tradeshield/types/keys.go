package types

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	// ModuleName defines the module name
	ModuleName = "tradeshield"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	ParamsKey                     = []byte{0x01}
	PendingSpotOrderKey           = []byte{0x02}
	PendingSpotOrderCountKey      = []byte{0x03}
	SortedSpotOrderKey            = []byte{0x04}
	PendingPerpetualOrderKey      = []byte{0x05}
	PendingPerpetualOrderCountKey = []byte{0x06}
	SortedPerpetualOrderKey       = []byte{0x07}
)

func GenSpotKey(order SpotOrder) (string, error) {
	if order.OrderType == SpotOrderType_MARKETBUY {
		return "", errors.New("cannot generate a key on a market order")
	}
	return fmt.Sprintf("%s\n%s\n%s",
		order.OrderType,
		order.OrderAmount.Denom,
		order.OrderTargetDenom), nil
}

func GenPerpKey(order PerpetualOrder) string {
	return fmt.Sprintf("%s\n%s\n%s",
		order.Position,
		order.PerpetualOrderType,
		order.TradingAsset)
}

func EncodeUint64Slice(slice []uint64) []byte {
	buf := make([]byte, 8*len(slice))
	for i, v := range slice {
		binary.BigEndian.PutUint64(buf[i*8:], v)
	}
	return buf
}

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

// GetUint64Bytes returns the byte representation of the ID
func GetUint64Bytes(ID uint64) []byte {
	IDBz := make([]byte, 8)
	binary.BigEndian.PutUint64(IDBz, ID)
	return IDBz
}

// GetUint64FromBytes returns ID in uint64 format from a byte array
func GetUint64FromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
