package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) OpenConsolidateLong(ctx sdk.Context, poolId uint64, mtp *types.MTP, msg *types.MsgOpen) (*types.MTP, error) {
	maxLeverage := k.OpenLongChecker.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)
	eta := leverage.Sub(sdk.OneDec())
	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	mtp.Leverage = leverage

	return k.ProcessOpenLong(ctx, mtp, leverage, eta, collateralAmountDec, poolId, msg)
}
