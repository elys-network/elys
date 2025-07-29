package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

func (k Keeper) OpenConsolidateMergeMtp(ctx sdk.Context, existingMtp *types.MTP, newMtp *types.MTP) (*types.MTP, error) {
	// Merge MTPs
	existingMtp.Collateral = existingMtp.Collateral.Add(newMtp.Collateral)
	existingMtp.Custody = existingMtp.Custody.Add(newMtp.Custody)
	existingMtp.Liabilities = existingMtp.Liabilities.Add(newMtp.Liabilities)

	// Destroy new MTP
	k.DestroyMTP(ctx, *newMtp)

	return existingMtp, nil
}
