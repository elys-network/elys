package types

const (
	// ModuleName defines the module name
	ModuleName = "oracle"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_oracle"

	// Version defines the current version the IBC module supports
	Version = "bandchain-1"

	// PortID is the default port id that module binds to
	PortID = "oracle"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("oracle-port-")
	// AssetInfoKeyPrefix is the prefix to retrieve all AssetInfo
	AssetInfoKeyPrefix = "AssetInfo/value/"
	// PriceKeyPrefix is the prefix to retrieve all Price
	PriceKeyPrefix = "Price/value/"
	// PriceFeederKeyPrefix is the prefix to retrieve all PriceFeeder
	PriceFeederKeyPrefix = "PriceFeeder/value/"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// AssetInfoKey returns the store key to retrieve a AssetInfo from the index fields
func AssetInfoKey(denom string) []byte {
	var key []byte

	indexBytes := []byte(denom)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}

// PriceKey returns the store key to retrieve a Price from the index fields
func PriceKey(asset string) []byte {
	var key []byte

	indexBytes := []byte(asset)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}

// PriceFeederKey returns the store key to retrieve a PriceFeeder from the feeder fields
func PriceFeederKey(feeder string) []byte {
	var key []byte

	indexBytes := []byte(feeder)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
