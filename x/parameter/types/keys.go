package types

import "encoding/binary"

const (
	// ModuleName defines the module name
	ModuleName = "parameter"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_parameter"

	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = "Params/value/"
)

const (
	// Eden denom
	Elys = "uelys"

	// Eden denom
	Eden = "ueden"

	// Eden Boost denom
	EdenB = "uedenb"

	// Base currency
	BaseCurrency = "uusdc"

	// Atom denom
	ATOM = "uatom"

	// BaseDecimal
	BASE_DECIMAL = 6

	// 3.154e+7 per year
	SecondsPerYear = 3.154e+7

	// One day seconds
	SecondsPerDay = 86400

	// Days per year
	DaysPerYear = 365

	// Weeks per year
	WeeksPerYear = 52

	// Days per week
	DaysPerWeek = 7

	// Return ok
	RES_OK = uint64(200)
)

var _ binary.ByteOrder

func KeyPrefix(p string) []byte {
	return []byte(p)
}
