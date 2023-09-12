package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) OpenLong(ctx sdk.Context, poolId uint64, msg *types.MsgOpen) (*types.MTP, error) {
	maxLeverage := k.OpenLongChecker.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)
	eta := leverage.Sub(sdk.OneDec())
	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	mtp := types.NewMTP(msg.Creator, msg.CollateralAsset, msg.BorrowAsset, msg.Position, leverage, poolId)

	return k.ProcessOpenLong(ctx, mtp, leverage, eta, collateralAmountDec, poolId, msg)
}
