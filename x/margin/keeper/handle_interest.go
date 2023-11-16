package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) HandleInterest(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, collateralAsset string, custodyAsset string) error {
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)
	if epochPosition <= 0 {
		return nil
	}

	interestPayment, err := k.CalcMTPInterestLiabilities(ctx, mtp, pool.InterestRate, epochPosition, epochLength, ammPool, collateralAsset)
	if err != nil {
		return err
	}
	finalInterestPayment := k.HandleInterestPayment(ctx, collateralAsset, custodyAsset, interestPayment, mtp, pool, ammPool)

	// finalInterestPayment is in custodyAsset
	if err := pool.UpdateBlockInterest(ctx, custodyAsset, finalInterestPayment, true, mtp.Position); err != nil {
		return err
	}

	_, err = k.UpdateMTPHealth(ctx, *mtp, ammPool)
	return err
}
