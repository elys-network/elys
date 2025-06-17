package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V12Migration(ctx sdk.Context) error {

	err := m.keeper.BurnAirdropWallet(ctx)
	if err != nil {
		// log error
		m.keeper.Logger(ctx).Error("failed to burn airdrop wallet", "error", err)
	}

	return nil
}
