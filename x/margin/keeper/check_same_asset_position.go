package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CheckSameAssetPosition(ctx sdk.Context, msg *types.MsgOpen) *types.MTP {
	mtps := k.GetAllMTPs(ctx)
	for _, mtp := range mtps {
		if mtp.Address == msg.Creator && mtp.Position == msg.Position && mtp.CollateralAsset == msg.Collateral.Denom && mtp.CustodyAsset == msg.TradingAsset {
			return &mtp
		}
	}

	return nil
}
