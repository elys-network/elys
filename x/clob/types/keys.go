package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/elys-network/elys/x/clob/utils"
	"math"
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

var (
	ParamsPrefix           = []byte{0x01}
	SubAccountPrefix       = []byte{0x02}
	PerpetualMarketPrefix  = []byte{0x03}
	PerpetualPrefix        = []byte{0x04}
	PerpetualOwnerPrefix   = []byte{0x05}
	PerpetualOrderPrefix   = []byte{0x06}
	MarketPricePrefix      = []byte{0x07}
	PerpetualCounterPrefix = []byte{0x08}
)

func GetAddressSubAccountPrefixKey(addr sdk.AccAddress) []byte {
	key := append(SubAccountPrefix, address.MustLengthPrefix(addr.Bytes())...)
	key = append(key, []byte("/")...)
	return key
}

func GetSubAccountKey(addr sdk.AccAddress, id uint64) []byte {
	key := GetAddressSubAccountPrefixKey(addr)
	key = append(key, sdk.Uint64ToBigEndian(id)...)
	return key
}

func GetPerpetualMarketKey(id uint64) []byte {
	return append(PerpetualMarketPrefix, sdk.Uint64ToBigEndian(id)...)
}

func GetPerpetualKey(id uint64) []byte {
	return append(PerpetualPrefix, sdk.Uint64ToBigEndian(id)...)
}

func GetPerpetualOwnerKey(subAccountId uint64, addr sdk.AccAddress, marketId uint64) []byte {
	key := append(PerpetualOwnerPrefix, address.MustLengthPrefix(addr.Bytes())...)
	key = append(key, []byte("/")...)
	key = append(key, sdk.Uint64ToBigEndian(subAccountId)...)
	key = append(key, []byte("/")...)
	key = append(key, sdk.Uint64ToBigEndian(marketId)...)
	return key
}

func GetPerpetualOrderKey(marketId uint64, orderType OrderType, price sdkmath.LegacyDec, height uint64) []byte {
	key := append(PerpetualOrderPrefix, sdk.Uint64ToBigEndian(marketId)...)
	key = append(key, []byte("/")...)
	orderTypeByte := FalseByte
	heightBytes := sdk.Uint64ToBigEndian(height)
	if IsBuy(orderType) {
		orderTypeByte = TrueByte
		heightBytes = sdk.Uint64ToBigEndian(math.MaxUint64 - height) // Subtracting it so that in buy order book it's sorted by height as Reverse iterator will be used
	}
	key = append(key, orderTypeByte)
	key = append(key, []byte("/")...)
	paddedPrice := utils.GetPaddedPrice(price)
	key = append(key, []byte(paddedPrice)...)
	key = append(key, []byte("/")...)
	key = append(key, heightBytes...)
	return key
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

func GetMarketPriceKey(id uint64) []byte {
	return append(MarketPricePrefix, sdk.Uint64ToBigEndian(id)...)
}

func GetPerpetualCounterKey(marketId uint64) []byte {
	return append(PerpetualCounterPrefix, sdk.Uint64ToBigEndian(marketId)...)
}
