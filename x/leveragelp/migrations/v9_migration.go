package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V9Migration(ctx sdk.Context) error {
	positions := m.keeper.GetAllPositions(ctx)
	for _, position := range positions {
		m.keeper.SetPosition(ctx, &position, sdk.ZeroInt())
		m.keeper.DeleteLegacyPosition(ctx, position.Address, position.Id)
	}

	whitelistAddressStrings := m.keeper.GetAllLegacyWhitelistedAddress(ctx)
	for _, addressString := range whitelistAddressStrings {
		m.keeper.WhitelistAddress(ctx, sdk.MustAccAddressFromBech32(addressString))
		m.keeper.DeleteLegacyWhitelistedAddress(ctx, addressString)
	}

	return nil
}
