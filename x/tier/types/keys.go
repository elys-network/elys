package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "tier"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_tier"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	PortfolioKeyPrefix = []byte{0x01}
	ParamKeyPrefix     = []byte{0x02}
)

func PortfolioKey(
	asset string,
) []byte {
	var key []byte

	assetBytes := []byte(asset)
	key = append(key, assetBytes...)
	key = append(key, []byte("/")...)

	return key
}

func GetPortfolioKey(date string, addr sdk.AccAddress) []byte {
	key := PortfolioKeyPrefix

	key = append(key, []byte(date)...)
	key = append(key, []byte("/")...)
	key = append(key, address.MustLengthPrefix(addr)...)
	return key
}

func GetPortfolioByDateIteratorKey(date string) []byte {
	key := PortfolioKeyPrefix

	key = append(key, []byte(date)...)
	key = append(key, []byte("/")...)
	return key
}
