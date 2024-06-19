package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CheckSameAssetPosition(ctx sdk.Context, msg *types.MsgOpen) *types.MTP {
	mtps, _, err := k.GetMTPsForAddress(ctx, sdk.MustAccAddressFromBech32(msg.Creator), &query.PageRequest{})
	if err != nil {
		return nil
	}
	for _, mtp := range mtps {
		if mtp.Position == msg.Position && mtp.CollateralAsset == msg.Collateral.Denom && mtp.CustodyAsset == msg.TradingAsset {
			return mtp
		}
	}

	return nil
}
