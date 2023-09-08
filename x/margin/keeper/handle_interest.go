package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) HandleInterest(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) error {
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)
	if epochPosition <= 0 {
		return nil
	}

	interestPayment := k.CalcMTPInterestLiabilities(mtp, pool.InterestRate, epochPosition, epochLength)
	finalInterestPayment := k.HandleInterestPayment(ctx, interestPayment, mtp, pool, ammPool)

	if err := pool.UpdateBlockInterest(ctx, mtp.CollateralAsset, finalInterestPayment, true); err != nil {
		return err
	}

	_, err := k.UpdateMTPHealth(ctx, *mtp, ammPool)
	return err
}
