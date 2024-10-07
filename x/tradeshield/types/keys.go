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
	if order.OrderPrice != nil {
		return fmt.Sprintf("%s\n%s\n%s",
			order.OrderType,
			order.OrderPrice.BaseDenom,
			order.OrderPrice.QuoteDenom), nil
	}
	return "", errors.New("order price not found")
}

func GenPerpKey(order PerpetualOrder) (string, error) {
	if order.PerpetualOrderType == PerpetualOrderType_MARKETCLOSE || order.PerpetualOrderType == PerpetualOrderType_MARKETOPEN {
		return "", errors.New("cannot generate a key on a market order")
	}
	if order.TriggerPrice != nil {
		return fmt.Sprintf("%s\n%s\n%s\n%s",
			order.Position,
			order.PerpetualOrderType,
			order.TriggerPrice.BaseDenom,
			order.TriggerPrice.QuoteDenom), nil
	}
	return "", errors.New("trigger price not found")
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
