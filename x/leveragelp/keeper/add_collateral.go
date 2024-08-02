package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

// Increase collateral, repay with additional collateral, update debt, liability and health
func (k Keeper) ProcessAddCollateral(ctx sdk.Context, address string, id uint64, collateral sdk.Int) error {
	position, err := k.GetPosition(ctx, address, id)
	if err != nil {
		return err
	}
	// Fetch the pool associated with the given pool ID.
	pool, found := k.GetPool(ctx, position.AmmPoolId)
	if !found {
		return errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", position.AmmPoolId))
	}

	// Check if the pool is enabled.
	if !k.IsPoolEnabled(ctx, position.AmmPoolId) {
		return errorsmod.Wrap(types.ErrPositionDisabled, fmt.Sprintf("poolId: %d", position.AmmPoolId))
	}

	oldDebt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress())
	// Check if collateral is not more than borrowed
	debtBefore := k.stableKeeper.UpdateInterestStackedByAddress(ctx, position.GetPositionAddress())
	maxAllowedCollateral := debtBefore.Borrowed.Add(debtBefore.InterestStacked).Sub(debtBefore.InterestPaid)
	if collateral.GT(maxAllowedCollateral) {
		return errorsmod.Wrap(types.ErrInvalidCollateral, fmt.Sprintf("Cannot add more than: %s", maxAllowedCollateral.String()))
	}

	// send collateral coins to Position address from Position owner address
	positionOwner := sdk.MustAccAddressFromBech32(position.Address)
	err = k.bankKeeper.SendCoins(ctx, positionOwner, position.GetPositionAddress(), sdk.Coins{sdk.NewCoin(position.Collateral.Denom, collateral)})
	if err != nil {
		return err
	}

	err = k.stableKeeper.Repay(ctx, position.GetPositionAddress(), sdk.NewCoin(position.Collateral.Denom, collateral))
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
	position.PositionHealth = positionHealth

	// Update Liabilities
	debt := k.stableKeeper.UpdateInterestStackedByAddress(ctx, position.GetPositionAddress())
	position.Liabilities = debt.Borrowed
	position.Collateral = position.Collateral.Add(sdk.NewCoin(position.Collateral.Denom, collateral))

	k.SetPosition(ctx, &position, oldDebt.Borrowed.Add(oldDebt.InterestStacked).Sub(oldDebt.InterestPaid))

	return nil
}
