package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetZeroAddress returns the zero address as a string
func GetZeroAddress() sdk.AccAddress {
	zeroBytes := make([]byte, 20)
	return sdk.AccAddress(zeroBytes)
}
