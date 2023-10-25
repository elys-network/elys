package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Get MTP address
func (p *MTP) GetMTPAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("leveragelp/%d", p.Id))
}
