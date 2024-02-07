package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OpenConsolidateLong(ctx sdk.Context, poolId uint64, existingMtp *types.MTP, newMtp *types.MTP, msg *types.MsgOpen, baseCurrency string) (*types.MTP, error) {
	existingMtp.Collateral = existingMtp.Collateral.Add(newMtp.Collateral)
	existingMtp.Custody = existingMtp.Custody.Add(newMtp.Custody)
	existingMtp.Liabilities = existingMtp.Liabilities.Add(newMtp.Liabilities)

	// Set existing MTP
	if err := k.SetMTP(ctx, existingMtp); err != nil {
		return nil, err
	}

	// Destroy new MTP
	if err := k.DestroyMTP(ctx, newMtp.Address, newMtp.Id); err != nil {
		return nil, err
	}

	return existingMtp, nil
}
