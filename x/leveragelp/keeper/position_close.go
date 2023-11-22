package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) ForceCloseLong(ctx sdk.Context, position types.Position, pool types.Pool) (sdk.Int, error) {
	// Exit liquidity with collateral token
	exitCoins, err := k.amm.ExitPool(ctx, position.GetPositionAddress(), position.AmmPoolId, position.LeveragedLpAmount, sdk.Coins{}, position.Collateral.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	// Repay with interest
	debt := k.stableKeeper.UpdateInterestStackedByAddress(ctx, position.GetPositionAddress())
	repayAmount := debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid)
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

	err = k.DestroyPosition(ctx, position.Address, position.Id)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	pool.LeveragedLpAmount = pool.LeveragedLpAmount.Sub(position.LeveragedLpAmount)
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
		return nil, sdk.ZeroInt(), sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
	}

	repayAmount, err := k.ForceCloseLong(ctx, position, pool)
	return &position, repayAmount, err
}
