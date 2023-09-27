package types

const (
	// ModuleName defines the module name
	ModuleName = "commitment"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_commitment"
)

const (
	BaseDenom       = "base-denom"
	VestingDenom    = "vesting-denom"
	EpochIdentifier = "epoch-identifier"
	NumEpochs       = "num-epochs"
	VestNowFactor   = "vest-now-factor"
	NumMaxVestings  = "num-max-vestings"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
