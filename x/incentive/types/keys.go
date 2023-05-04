package types

const (
	// ModuleName defines the module name
	ModuleName = "incentive"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName + "_store"

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_incentive"
)

// Denom definition - Eden and Eden boost
const (
	// Eden denom
	Eden = "ueden"

	// Eden Boost denom
	EdenB = "uedenb"

	// 52.1429 weeks
	WeeksPerYear = 52

	// 365 days
	DaysPerYear = 365

	// 8760 hours
	HoursPerYear = 8760
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
