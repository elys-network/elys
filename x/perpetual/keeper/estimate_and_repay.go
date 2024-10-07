package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) EstimateAndRepay(ctx sdk.Context, mtp types.MTP, pool types.Pool, ammPool ammtypes.Pool, amount math.Int, baseCurrency string) (math.Int, error) {
	// init repay amount
	repayAmount := sdk.ZeroInt()
	var err error

	// if position is long, repay in collateral asset
	if amount.IsPositive() {
		if mtp.Position == types.Position_LONG {
			custodyAmtTokenIn := sdk.NewCoin(mtp.CustodyAsset, amount)
			repayAmount, err = k.EstimateSwap(ctx, custodyAmtTokenIn, mtp.CollateralAsset, ammPool)
			if err != nil {
				return sdk.ZeroInt(), err
			}
		} else if mtp.Position == types.Position_SHORT {
			// if position is short, repay in trading asset
			custodyAmtTokenIn := sdk.NewCoin(mtp.CustodyAsset, amount)
			repayAmount, err = k.EstimateSwap(ctx, custodyAmtTokenIn, mtp.TradingAsset, ammPool)
			if err != nil {
				return sdk.ZeroInt(), err
			}
		} else {
			return sdk.ZeroInt(), types.ErrInvalidPosition
		}
	}

	returnAmount, err := k.CalcReturnAmount(ctx, mtp, pool, ammPool, repayAmount, amount, baseCurrency)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	// update mtp health
	mtp.MtpHealth, err = k.GetMTPHealth(ctx, mtp, ammPool, baseCurrency)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	// if return amount positive then update liabilities
	if returnAmount.IsPositive() {
		mtp.Liabilities = mtp.Liabilities.Sub(mtp.Liabilities.Mul(amount).Quo(mtp.Custody))
	}

	// Note: Long settlement is done in trading asset. And short settlement in usdc in Repay function
	if err := k.Repay(ctx, &mtp, &pool, ammPool, returnAmount, amount); err != nil {
		return sdk.ZeroInt(), err
	}

	return repayAmount, nil
}
