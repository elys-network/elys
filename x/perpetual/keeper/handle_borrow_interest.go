package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) HandleBorrowInterest(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) error {
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)
	if epochPosition <= 0 {
		return nil
	}

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	borrowInterestPayment, err := k.CalcMTPBorrowInterestLiabilities(ctx, mtp, pool.BorrowInterestRate, epochPosition, epochLength, ammPool, baseCurrency)
	if err != nil {
		return err
	}
	finalBorrowInterestPayment := k.HandleBorrowInterestPayment(ctx, borrowInterestPayment, mtp, pool, ammPool, baseCurrency)

	// finalInterestPayment is in custodyAsset
	if err := pool.UpdateBlockBorrowInterest(ctx, mtp.CustodyAsset, finalBorrowInterestPayment, true, mtp.Position); err != nil {
		return err
	}

	_, err = k.UpdateMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	return err
}
