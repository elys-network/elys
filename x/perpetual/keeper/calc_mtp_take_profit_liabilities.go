package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CalcMTPTakeProfitLiability(ctx sdk.Context, mtp *types.MTP, baseCurrency string) (math.Int, error) {
	if mtp.TakeProfitCustody.IsZero() {
		return math.ZeroInt(), nil
	}

	// Retrieve AmmPool
	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	// build take profit custody coin
	takeProfitCustody := sdk.NewCoin(mtp.CustodyAsset, mtp.TakeProfitCustody)

	// convert custody amount to base currency
	takeProfitLiabilities, _, err := k.EstimateSwap(ctx, takeProfitCustody, baseCurrency, ammPool)
	if err != nil {
		return sdk.ZeroInt(), errorsmod.Wrapf(err, "unable to swap takeProfitCustody to baseCurrency for takeProfitLiabilities")
	}

	return takeProfitLiabilities, nil
}
