package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func GetSpotOrderAddress(orderId uint64) sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("tradeshield/spot/%d", orderId))
}

func GetPerpOrderAddress(orderId uint64) sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("tradeshield/perp/%d", orderId))
}

// Get Spot Order address
func (p SpotOrder) GetOrderAddress() sdk.AccAddress {
	return GetSpotOrderAddress(p.OrderId)
}

// Get Perpetual Order address
func (p PerpetualOrder) GetOrderAddress() sdk.AccAddress {
	return GetPerpOrderAddress(p.OrderId)
}
