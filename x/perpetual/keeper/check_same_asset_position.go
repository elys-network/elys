package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
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

func (k Keeper) CheckSamePositionAndConsolidate(ctx sdk.Context, m *types.MTP) error {
	mtps := k.GetAllMTPsForAddress(ctx, sdk.AccAddress(m.Address))
	for _, mtp := range mtps {
		if mtp.Position == m.Position && mtp.CollateralAsset == m.CollateralAsset && mtp.CustodyAsset == m.TradingAsset {
			if mtp.Id != m.Id {

				_, err := k.OpenConsolidateMergeMtp(ctx, m.AmmPoolId, mtp, m, ptypes.BaseCurrency)
				if err != nil {
					return err
				}

				ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId, mtp.CustodyAsset)
				if err != nil {
					return err
				}
				entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
				if !found {
					return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
				}
				baseCurrency := entry.Denom
				// calc and update open price
				err = k.UpdateOpenPrice(ctx, mtp, ammPool, baseCurrency)
				if err != nil {
					return err
				}
				return nil
			}

		}
	}
	return nil
}
