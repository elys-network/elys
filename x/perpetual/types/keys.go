package types

import "encoding/binary"

const (
	// ModuleName defines the module name
	ModuleName = "perpetual"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_perpetual"

	// ParamsKey is the prefix for parameters of perpetual module
	ParamsKey = "perpetual_params"
)

const MaxPageLimit = 100

const (
	InfinitePriceString    = "infinite"
	TakeProfitPriceDefault = "10000000000000000000000000000000000000000" // 10^40
)

var (
	MTPPrefix          = []byte{0x01}
	MTPCountPrefix     = []byte{0x02}
	OpenMTPCountPrefix = []byte{0x04}
	WhitelistPrefix    = []byte{0x05}
	InterestRatePrefix = []byte{0x06}
	FundingRatePrefix  = []byte{0x07}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
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

func GetWhitelistKey(address string) []byte {
	return append(WhitelistPrefix, []byte(address)...)
}

func GetMTPKey(address string, id uint64) []byte {
	return append(MTPPrefix, append([]byte(address), GetUint64Bytes(id)...)...)
}

func GetMTPPrefixForAddress(address string) []byte {
	return append(MTPPrefix, []byte(address)...)
}

func GetInterestRateKey(block uint64, pool uint64) []byte {
	return append(GetUint64Bytes(block), GetUint64Bytes(pool)...)
}

func GetFundingRateKey(block uint64, pool uint64) []byte {
	return append(GetUint64Bytes(block), GetUint64Bytes(pool)...)
}
