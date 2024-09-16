package types

import (
	"errors"
	"fmt"
)

const (
	// ModuleName defines the module name
	ModuleName = "tradeshield"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_tradeshield"
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
