package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CalcMTPConsolidateCollateral(ctx sdk.Context, mtp *types.MTP, baseCurrency string) error {
	consolidateCollateral := sdk.ZeroInt()
	if mtp.CollateralAsset == baseCurrency {
		consolidateCollateral = consolidateCollateral.Add(mtp.Collateral)
	} else {
		// swap into base currency
		_, ammPool, _, err := k.OpenChecker.PreparePools(ctx, mtp.CollateralAsset, mtp.CustodyAsset)
		if err != nil {
			return err
		}

		collateralAmtIn := sdk.NewCoin(mtp.CollateralAsset, mtp.Collateral)
		C, err := k.EstimateSwapGivenOut(ctx, collateralAmtIn, baseCurrency, ammPool)
		if err != nil {
			return err
		}

		consolidateCollateral = consolidateCollateral.Add(C)
	}

	mtp.SumCollateral = consolidateCollateral

	return nil
}
