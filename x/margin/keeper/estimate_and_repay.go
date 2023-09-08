package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) EstimateAndRepay(ctx sdk.Context, mtp types.MTP, pool types.Pool, ammPool ammtypes.Pool) (sdk.Int, error) {
	cutodyAmtTokenIn := sdk.NewCoin(mtp.CustodyAsset, mtp.CustodyAmount)
	repayAmount, err := k.EstimateSwap(ctx, cutodyAmtTokenIn, mtp.CollateralAsset, ammPool)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	if err := k.Repay(ctx, &mtp, &pool, ammPool, repayAmount, false); err != nil {
		return sdk.ZeroInt(), err
	}

	return repayAmount, nil
}
