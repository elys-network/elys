package constants

import authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

var ZeroAddress = ""

// SetZeroAddress Should be called after setting prefixes for address
func SetZeroAddress() {
	ZeroAddress = authtypes.NewModuleAddress("zero").String()
}
