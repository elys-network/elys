package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (s SubAccount) GetOwnerAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(s.Owner)
}

func (s SubAccount) GetTradingAccountAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("clob/%s/%d", s.Owner, s.Id))
}
