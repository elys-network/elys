package types

var (
	// BandPriceResultStoreKeyPrefix is a prefix for storing result
	BandPriceResultStoreKeyPrefix = "band_price_result"

	// LastBandRequestIdKey is the key for the last request id
	LastBandRequestIdKey = "last_band_request_id"

	// BandPriceClientIDKey is query request identifier
	BandPriceClientIDKey = "band_price_id"

	// PrefixKeyBandRequest is the prefix for band requests
	PrefixKeyBandRequest = "band_request_"
)

// BandPriceResultStoreKey is a function to generate key for each result in store
func BandPriceResultStoreKey(requestID OracleRequestID) []byte {
	return append(KeyPrefix(BandPriceResultStoreKeyPrefix), int64ToBytes(int64(requestID))...)
}

func BandRequestStoreKey(requestID OracleRequestID) []byte {
	return append(KeyPrefix(PrefixKeyBandRequest), int64ToBytes(int64(requestID))...)
}
