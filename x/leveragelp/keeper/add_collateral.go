package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/leveragelp/types"
)

// Increase collateral, repay with additional collateral, update debt, liability and health
func (k Keeper) ProcessAddCollateral(ctx sdk.Context, address string, id uint64, collateral sdkmath.Int) error {
	creator := sdk.MustAccAddressFromBech32(address)
	position, err := k.GetPosition(ctx, creator, id)
	if err != nil {
		return err
	}
	// Fetch the pool associated with the given pool ID.
	pool, found := k.GetPool(ctx, position.AmmPoolId)
	if !found {
		return errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", position.AmmPoolId))
	}

	// Check if collateral is not more than borrowed
	debtBefore := k.stableKeeper.UpdateInterestAndGetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId, position.AmmPoolId)
	maxAllowedCollateral := debtBefore.GetTotalLiablities()
	if collateral.GT(maxAllowedCollateral) {
		return errorsmod.Wrap(types.ErrInvalidCollateral, fmt.Sprintf("Cannot add more than: %s", maxAllowedCollateral.String()))
	}

	// send collateral coins to Position address from Position owner address
	err = k.bankKeeper.SendCoins(ctx, position.GetOwnerAddress(), position.GetPositionAddress(), sdk.Coins{sdk.NewCoin(position.Collateral.Denom, collateral)})
	if err != nil {
		return err
	}

	err = k.stableKeeper.Repay(ctx, position.GetPositionAddress(), sdk.NewCoin(position.Collateral.Denom, collateral), position.BorrowPoolId, position.AmmPoolId)
	if err != nil {
		return err
	}

	// Update the pool health.
	k.UpdatePoolHealth(ctx, &pool)

	// Update position health
	positionHealth, err := k.GetPositionHealth(ctx, position)
	if err != nil {
		return err
	}
	position.PositionHealth = positionHealth.Dec()

	// Update Liabilities
	debt := k.stableKeeper.UpdateInterestAndGetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId, position.AmmPoolId)
	position.Liabilities = debt.GetTotalLiablities()
	position.Collateral = position.Collateral.Add(sdk.NewCoin(position.Collateral.Denom, collateral))

	k.SetPosition(ctx, &position)

	return nil
}
