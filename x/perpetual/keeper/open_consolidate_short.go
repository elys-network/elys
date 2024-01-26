package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OpenConsolidateShort(ctx sdk.Context, poolId uint64, mtp *types.MTP, msg *types.MsgOpen, baseCurrency string) (*types.MTP, error) {
	if mtp.Leverage != msg.Leverage {
		return nil, types.ErrInvalidLeverage
	}

	maxLeverage := k.OpenLongChecker.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(mtp.Leverage, maxLeverage)
	eta := leverage.Sub(sdk.OneDec())
	collateralAmountDec := sdk.NewDecFromBigInt(msg.Collateral.Amount.BigInt())

	return k.ProcessOpenShort(ctx, mtp, leverage, eta, collateralAmountDec, poolId, msg, baseCurrency)
}
