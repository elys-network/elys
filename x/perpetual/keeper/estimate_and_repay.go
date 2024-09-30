package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) EstimateAndRepay(ctx sdk.Context, mtp types.MTP, pool types.Pool, ammPool ammtypes.Pool, amount sdkmath.Int, baseCurrency string) (sdkmath.Int, error) {
	// init repay amount
	var repayAmount sdkmath.Int
	var err error

	// if position is long, repay in collateral asset
	if mtp.Position == types.Position_LONG {
		custodyAmtTokenIn := sdk.NewCoin(mtp.CustodyAsset, amount)
		repayAmount, err = k.EstimateSwap(ctx, custodyAmtTokenIn, mtp.CollateralAsset, ammPool)
		if err != nil {
			return sdkmath.ZeroInt(), err
		}
	} else if mtp.Position == types.Position_SHORT {
		// if position is short, repay in trading asset
		custodyAmtTokenIn := sdk.NewCoin(mtp.CustodyAsset, amount)
		repayAmount, err = k.EstimateSwap(ctx, custodyAmtTokenIn, mtp.TradingAsset, ammPool)
		if err != nil {
			return sdkmath.ZeroInt(), err
		}
	} else {
		return sdkmath.ZeroInt(), types.ErrInvalidPosition
	}

	// Note: Long settlement is done in trading asset. And short settlement in usdc in Repay function
	if err := k.Repay(ctx, &mtp, &pool, ammPool, repayAmount, false, amount, baseCurrency); err != nil {
		return sdkmath.ZeroInt(), err
	}

	return repayAmount, nil
}
