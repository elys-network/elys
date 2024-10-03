package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CheckSameAssetPosition(ctx sdk.Context, msg *types.MsgOpen) *types.MTP {
	mtps := k.GetAllMTPsForAddress(ctx, sdk.MustAccAddressFromBech32(msg.Creator))
	for _, mtp := range mtps {
		if mtp.Position == msg.Position && mtp.CollateralAsset == msg.Collateral.Denom && mtp.CustodyAsset == msg.TradingAsset {
			return mtp
		}
	}

	return nil
}

func (k Keeper) CheckSamePositionAndRemove(ctx sdk.Context, m *types.MTP) error {
	mtps := k.GetAllMTPsForAddress(ctx, sdk.AccAddress(m.Address))
	for _, mtp := range mtps {
		if mtp.Position == m.Position && mtp.CollateralAsset == m.CollateralAsset && mtp.CustodyAsset == m.TradingAsset {
			if mtp.Id != m.Id {
				switch m.Position {
				case types.Position_LONG:
					_, err := k.OpenConsolidateLong(ctx, m.AmmPoolId, mtp, m)
					if err != nil {
						return err
					}
				case types.Position_SHORT:
					_, err := k.OpenConsolidateShort(ctx, m.AmmPoolId, mtp, m)
					if err != nil {
						return err
					}
				}
				
			}

		}
	}
	return nil
}
