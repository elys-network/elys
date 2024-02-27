package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) UpdateOpenPrice(ctx sdk.Context, mtp *types.MTP, ammPool ammtypes.Pool, baseCurrency string) error {
	collateralAmountInBaseCurrency := mtp.Collateral
	if mtp.CollateralAsset != baseCurrency {
		C, err := k.EstimateSwap(ctx, sdk.NewCoin(mtp.CollateralAsset, mtp.Collateral), baseCurrency, ammPool)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error estimating swap: %s", mtp.CustodyAsset))
		}
		collateralAmountInBaseCurrency = C
	}

	// open price = (collateral + liabilities) / custody
	mtp.OpenPrice = math.LegacyNewDecFromBigInt(
		collateralAmountInBaseCurrency.Add(mtp.Liabilities).Quo(mtp.Custody).BigInt(),
	)

	err := k.SetMTP(ctx, mtp)
	if err != nil {
		return err
	}

	return nil
}
