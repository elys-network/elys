package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) CalcMTPConsolidateCollateral(ctx sdk.Context, mtp *types.MTP) error {
	consolidateCollateral := sdk.ZeroInt()
	for i, asset := range mtp.CollateralAssets {
		if asset == ptypes.BaseCurrency {
			consolidateCollateral = consolidateCollateral.Add(mtp.CollateralAmounts[i])
		} else {
			ammPool, found := k.amm.GetPool(ctx, mtp.AmmPoolId)
			if !found {
				return types.ErrAmmPoolNotFound.Wrap(fmt.Sprintf("poolId: %d", mtp.AmmPoolId))
			}

			collateralAmtIn := sdk.NewCoin(asset, mtp.CollateralAmounts[i])
			C, err := k.EstimateSwapGivenOut(ctx, collateralAmtIn, ptypes.BaseCurrency, ammPool)
			if err != nil {
				return err
			}

			consolidateCollateral = consolidateCollateral.Add(C)
		}
	}

	mtp.SumCollateral = consolidateCollateral

	return nil
}
