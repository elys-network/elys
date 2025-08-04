package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "clob"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	TStoreKey = "transient_" + ModuleName

	TrueByte  = byte(1)
	FalseByte = byte(0)
)

const MaxPageLimit = 10000

var (
	ParamsPrefix                 = []byte{0x01}
	SubAccountPrefix             = []byte{0x02}
	PerpetualMarketPrefix        = []byte{0x03}
	PerpetualMarketCounterPrefix = []byte{0x04}
	PerpetualPrefix              = []byte{0x05}
	PerpetualOwnerPrefix         = []byte{0x06}
	PerpetualOrderPrefix         = []byte{0x07}
	TwapPricesPrefix             = []byte{0x08}
	FundingRatePrefix            = []byte{0x09}
	PerpetualADLPrefix           = []byte{0x10}
	PerpetualOrderOwnerPrefix    = []byte{0x11}
)

func GetAddressSubAccountPrefixKey(addr sdk.AccAddress) []byte {
	key := append(SubAccountPrefix, address.MustLengthPrefix(addr.Bytes())...)
	key = append(key, []byte("/")...)
	return key
}

func GetSubAccountKey(addr sdk.AccAddress, subAccountId uint64) []byte {
	key := GetAddressSubAccountPrefixKey(addr)
	key = append(key, sdk.Uint64ToBigEndian(subAccountId)...)
	return key
}

func GetPerpetualMarketKey(id uint64) []byte {
	return append(PerpetualMarketPrefix, sdk.Uint64ToBigEndian(id)...)
}

func GetPerpetualMarketCounterKey(id uint64) []byte {
	return append(PerpetualMarketCounterPrefix, sdk.Uint64ToBigEndian(id)...)
}

func GetPerpetualKey(marketId, id uint64) []byte {
	key := append(PerpetualPrefix, sdk.Uint64ToBigEndian(marketId)...)
	key = append(key, []byte("/")...)
	return append(key, sdk.Uint64ToBigEndian(id)...)
}

func GetPerpetualOwnerAddressKey(addr sdk.AccAddress) []byte {
	key := append(PerpetualOwnerPrefix, address.MustLengthPrefix(addr.Bytes())...)
	key = append(key, []byte("/")...)
	return key
}

func GetPerpetualOwnerSubAccountKey(addr sdk.AccAddress, subAccountId uint64) []byte {
	key := GetPerpetualOwnerAddressKey(addr)
	key = append(key, sdk.Uint64ToBigEndian(subAccountId)...)
	key = append(key, []byte("/")...)
	return key
}

func GetPerpetualOwnerKey(addr sdk.AccAddress, subAccountId, marketId, perpetualId uint64) []byte {
	key := GetPerpetualOwnerSubAccountKey(addr, subAccountId)
	key = append(key, sdk.Uint64ToBigEndian(marketId)...)
	key = append(key, []byte("/")...)
	key = append(key, sdk.Uint64ToBigEndian(perpetualId)...)
	return key
}

func GetPerpetualOrderKey(marketId uint64, orderType OrderType, priceTick int64, counter uint64) []byte {
	return append(PerpetualOrderPrefix, NewOrderKey(marketId, orderType, priceTick, counter).KeyWithoutPrefix()...)
}

func GetPerpetualOrderBookIteratorKey(marketId uint64, long bool) []byte {
	key := PerpetualOrderPrefix
	key = append(key, sdk.Uint64ToBigEndian(marketId)...)
	key = append(key, []byte("/")...)
	orderTypeByte := FalseByte
	if long {
		orderTypeByte = TrueByte
	}
	key = append(key, orderTypeByte)
	key = append(key, []byte("/")...)
	return key
}

func GetFundingRateKey(id uint64) []byte {
	return append(FundingRatePrefix, sdk.Uint64ToBigEndian(id)...)
}

func GetTwapPricesKey(marketId, block uint64) []byte {
	key := append(TwapPricesPrefix, sdk.Uint64ToBigEndian(marketId)...)
	key = append(key, []byte("/")...)
	return append(key, sdk.Uint64ToBigEndian(block)...)
}

func GetPerpetualADLKey(marketId, id uint64) []byte {
	key := append(PerpetualADLPrefix, sdk.Uint64ToBigEndian(marketId)...)
	key = append(key, []byte("/")...)
	return append(key, sdk.Uint64ToBigEndian(id)...)
}

func GetOrderOwnerAddressKey(addr sdk.AccAddress) []byte {
	key := append(PerpetualOrderOwnerPrefix, address.MustLengthPrefix(addr.Bytes())...)
	key = append(key, []byte("/")...)
	return key
}

func GetOrderSubAccountKey(addr sdk.AccAddress, subAccountId uint64) []byte {
	key := GetOrderOwnerAddressKey(addr)
	key = append(key, sdk.Uint64ToBigEndian(subAccountId)...)
	key = append(key, []byte("/")...)
	return key
}

func GetOrderOwnerKey(addr sdk.AccAddress, subAccountId uint64, orderKey OrderKey) []byte {
	key := GetOrderSubAccountKey(addr, subAccountId)
	key = append(key, orderKey.KeyWithoutPrefix()...)
	return key
}
