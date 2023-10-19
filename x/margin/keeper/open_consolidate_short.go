package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) OpenConsolidateShort(ctx sdk.Context, poolId uint64, mtp *types.MTP, msg *types.MsgOpen) (*types.MTP, error) {
	maxLeverage := k.OpenShortChecker.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)
	eta := leverage.Sub(sdk.OneDec())
	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	mtp.Leverages = append(mtp.Leverages, leverage)

	mtp.Collaterals = types.AddOrAppendCoin(mtp.Collaterals, sdk.NewCoin(msg.CollateralAsset, sdk.NewInt(0)))
	mtp.Custodies = types.AddOrAppendCoin(mtp.Custodies, sdk.NewCoin(msg.BorrowAsset, sdk.NewInt(0)))

	return k.ProcessOpenShort(ctx, mtp, leverage, eta, collateralAmountDec, poolId, msg)
}
