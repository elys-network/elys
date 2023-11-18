package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) EstimateAndRepay(ctx sdk.Context, mtp types.MTP, pool types.Pool, ammPool ammtypes.Pool, collateralAsset string, custodyAsset string) (sdk.Int, error) {
	collateralIndex, custodyIndex := types.GetMTPAssetIndex(&mtp, collateralAsset, custodyAsset)
	cutodyAmtTokenIn := sdk.NewCoin(mtp.Custodies[custodyIndex].Denom, mtp.Custodies[custodyIndex].Amount)
	repayAmount, err := k.EstimateSwap(ctx, cutodyAmtTokenIn, mtp.Collaterals[collateralIndex].Denom, ammPool)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	if err := k.Repay(ctx, &mtp, &pool, ammPool, repayAmount, false, collateralAsset); err != nil {
		return sdk.ZeroInt(), err
	}

	return repayAmount, nil
}
