package types

import "encoding/binary"

const (
	// ModuleName defines the module name
	ModuleName = "parameter"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
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

	// Return ok
	RES_OK = uint64(200)

	USDC_DISPLAY = "USDC"
	USDT_DISPLAY = "USDT"
)

var _ binary.ByteOrder

func KeyPrefix(p string) []byte {
	return []byte(p)
}
