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

	// AnteHandlerParamKeyPrefix is the prefix to retrieve all AnteHandlerParam
	AnteHandlerParamKeyPrefix = "AnteHandlerParam/value/"

	AnteStoreKey = "ante-handler-param"
)

const (
	// Eden denom
	Elys = "uelys"

	// Eden denom
	Eden = "ueden"

	// Eden Boost denom
	EdenB = "uedenb"

	// USDC
	USDC = "cusdc"

	// 52.1429 weeks
	WeeksPerYear = 52

	// 365 days
	DaysPerYear = 365

	// 8760 hours
	HoursPerYear = 8760
)

var _ binary.ByteOrder

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// AnteHandlerParamKey returns the store key to retrieve a AnteHandlerParam from the index fields
func AnteHandlerParamKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
