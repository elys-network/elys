package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) EstimateAndRepay(ctx sdk.Context, mtp types.MTP, pool types.Pool, ammPool ammtypes.Pool, amount sdk.Int, baseCurrency string) (sdk.Int, error) {
	if amount.GT(mtp.Custody) || amount.IsNegative() {
		return sdk.ZeroInt(), types.ErrInvalidCloseSize
	}

	// init repay amount
	repayAmount := sdk.ZeroInt()

	// if position is long, repay in collateral asset
	if mtp.Position == types.Position_LONG {
		custodyAmtTokenIn := sdk.NewCoin(mtp.CustodyAsset, amount)
		var err error
		repayAmount, err = k.EstimateSwap(ctx, custodyAmtTokenIn, mtp.CollateralAsset, ammPool)
		if err != nil {
			return sdk.ZeroInt(), err
		}
	} else if mtp.Position == types.Position_SHORT {
		// if position is short, repay in trading asset
		custodyAmtTokenIn := sdk.NewCoin(mtp.CustodyAsset, amount)
		var err error
		repayAmount, err = k.EstimateSwap(ctx, custodyAmtTokenIn, mtp.TradingAsset, ammPool)
		if err != nil {
			return sdk.ZeroInt(), err
		}
	} else {
		return sdk.ZeroInt(), types.ErrInvalidPosition
	}

	if err := k.Repay(ctx, &mtp, &pool, ammPool, repayAmount, false, amount, baseCurrency); err != nil {
		return sdk.ZeroInt(), err
	}

	return repayAmount, nil
}
