package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CreateModuleAccount creates a module account at the provided address.
// It overrides an account if it exists at that address, with a non-zero sequence number & pubkey
// Contract: addr is derived from `address.Module(ModuleName, key)`
func CreateModuleAccount(ctx sdk.Context, ak types.AccountKeeper, addr sdk.AccAddress) error {
	err := CanCreateModuleAccountAtAddr(ctx, ak, addr)
	if err != nil {
		return err
	}

	acc := ak.NewAccount(
		ctx,
		authtypes.NewModuleAccount(
			authtypes.NewBaseAccountWithAddress(addr),
			addr.String(),
		),
	)
	ak.SetAccount(ctx, acc)
	return nil
}
