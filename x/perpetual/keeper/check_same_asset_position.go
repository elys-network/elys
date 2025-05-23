package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/perpetual/types"
)

func (k Keeper) CheckSameAssetPosition(ctx sdk.Context, msg *types.MsgOpen) *types.MTP {
	mtps := k.GetAllMTPsForAddress(ctx, sdk.MustAccAddressFromBech32(msg.Creator))
	for _, mtp := range mtps {
		if mtp.Position == msg.Position && mtp.CollateralAsset == msg.Collateral.Denom && mtp.TradingAsset == msg.TradingAsset {
			return mtp
		}
	}
	return nil
}
