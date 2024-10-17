package constants

import authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

var ZeroAddress = authtypes.NewModuleAddress("zero").String() // prefix will be in cosmos1... so this will fail, updated below

// SetZeroAddress Should be called after setting prefixes for address
func SetZeroAddress() {
	ZeroAddress = authtypes.NewModuleAddress("zero").String()
}
