package types

import sdkmath "cosmossdk.io/math"

type Trade struct {
	BuyerSubAccount, SellerSubAccount SubAccount
	MarketId                          uint64
	Price                             sdkmath.LegacyDec
	Quantity                          sdkmath.LegacyDec
}
