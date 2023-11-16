package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) HandleInterest(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, collateralAsset string, custodyAsset string) error {
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)
	if epochPosition <= 0 {
		return nil
	}

	entry, found := k.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	interestPayment := k.CalcMTPInterestLiabilities(ctx, mtp, pool.InterestRate, epochPosition, epochLength, ammPool, collateralAsset, baseCurrency)
	finalInterestPayment := k.HandleInterestPayment(ctx, collateralAsset, custodyAsset, interestPayment, mtp, pool, ammPool, baseCurrency)

	// finalInterestPayment is in custodyAsset
	if err := pool.UpdateBlockInterest(ctx, custodyAsset, finalInterestPayment, true, mtp.Position); err != nil {
		return err
	}

	_, err := k.UpdateMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	return err
}
