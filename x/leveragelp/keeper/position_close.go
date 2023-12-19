package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) ForceCloseLong(ctx sdk.Context, position types.Position, pool types.Pool, lpAmount sdk.Int) (sdk.Int, error) {
	if lpAmount.GT(position.LeveragedLpAmount) || lpAmount.IsNegative() {
		return sdk.ZeroInt(), types.ErrInvalidCloseSize
	}

	// Exit liquidity with collateral token
	exitCoins, err := k.amm.ExitPool(ctx, position.GetPositionAddress(), position.AmmPoolId, lpAmount, sdk.Coins{}, position.Collateral.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	// Repay with interest
	debt := k.stableKeeper.UpdateInterestStackedByAddress(ctx, position.GetPositionAddress())

	// Ensure position.LeveragedLpAmount is not zero to avoid division by zero
	if position.LeveragedLpAmount.IsZero() {
		return sdk.ZeroInt(), types.ErrAmountTooLow
	}

	repayAmount := debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid).Mul(lpAmount).Quo(position.LeveragedLpAmount)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	err = k.stableKeeper.Repay(ctx, position.GetPositionAddress(), sdk.NewCoin(position.Collateral.Denom, repayAmount))
	if err != nil {
		return sdk.ZeroInt(), err
	}

	userAmount := exitCoins[0].Amount.Sub(repayAmount)
	positionOwner := sdk.MustAccAddressFromBech32(position.Address)
	err = k.bankKeeper.SendCoins(ctx, positionOwner, position.GetPositionAddress(), sdk.Coins{sdk.NewCoin(position.Collateral.Denom, userAmount)})
	if err != nil {
		return sdk.ZeroInt(), err
	}

	position.LeveragedLpAmount = position.LeveragedLpAmount.Sub(lpAmount)
	if position.LeveragedLpAmount.IsZero() {
		err = k.DestroyPosition(ctx, position.Address, position.Id)
		if err != nil {
			return sdk.ZeroInt(), err
		}
	} else {
		k.SetPosition(ctx, &position)
	}

	pool.LeveragedLpAmount = pool.LeveragedLpAmount.Sub(lpAmount)
	k.UpdatePoolHealth(ctx, &pool)

	// Hooks after leveragelp position closed
	if k.hooks != nil {
		k.hooks.AfterLeveragelpPositionClosed(ctx, pool)
	}
	return repayAmount, nil
}

func (k Keeper) CloseLong(ctx sdk.Context, msg *types.MsgClose) (*types.Position, sdk.Int, error) {
	// Retrieve Position
	position, err := k.GetPosition(ctx, msg.Creator, msg.Id)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Retrieve Pool
	pool, found := k.GetPool(ctx, position.AmmPoolId)
	if !found {
		return nil, sdk.ZeroInt(), errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
	}

	// If lpAmount is lower than zero, close full amount
	lpAmount := msg.LpAmount
	if lpAmount.LTE(sdk.ZeroInt()) {
		lpAmount = position.LeveragedLpAmount
	}

	repayAmount, err := k.ForceCloseLong(ctx, position, pool, lpAmount)
	return &position, repayAmount, err
}
