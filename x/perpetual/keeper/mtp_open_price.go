package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) UpdateOpenPrice(ctx sdk.Context, mtp *types.MTP) error {
	mtp.GetAndSetOpenPrice()

	err := k.SetMTP(ctx, mtp)
	if err != nil {
		return err
	}

	return nil
}
