package keeper

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) OpenLong(ctx sdk.Context, msg *types.MsgOpen) (*types.Position, error) {
	// Determine the maximum leverage available and compute the effective leverage to be used.
	maxLeverage := k.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)

	// Convert the collateral amount into a decimal format.
	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())

	// Initialize a new Leveragelp Trading Position (Position).
	position := types.NewPosition(msg.Creator, sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount), leverage, msg.AmmPoolId)
	position.Id = k.GetPositionCount(ctx) + 1
	k.SetPositionCount(ctx, position.Id)

	// Call the function to process the open long logic.
	return k.ProcessOpenLong(ctx, position, leverage, collateralAmountDec, msg.AmmPoolId, msg)
}

func (k Keeper) OpenConsolidate(ctx sdk.Context, position *types.Position, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	poolId := position.AmmPoolId
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	if !k.IsPoolEnabled(ctx, poolId) {
		return nil, sdkerrors.Wrap(types.ErrPositionDisabled, fmt.Sprintf("poolId: %d", poolId))
	}

	ammPool, err := k.GetAmmPool(ctx, poolId)
	if err != nil {
		return nil, err
	}

	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	position.Collateral = position.Collateral.Add(sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount))
	maxLeverage := k.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)
	position.Leverage = leverage

	position, err = k.ProcessOpenLong(ctx, position, leverage, collateralAmountDec, poolId, msg)
	if err != nil {
		return nil, err
	}

	event := sdk.NewEvent(types.EventOpen,
		sdk.NewAttribute("id", strconv.FormatInt(int64(position.Id), 10)),
		sdk.NewAttribute("address", position.Address),
		sdk.NewAttribute("collateral", position.Collateral.String()),
		sdk.NewAttribute("leverage", position.Leverage.String()),
		sdk.NewAttribute("liabilities", position.Liabilities.String()),
		sdk.NewAttribute("health", position.PositionHealth.String()),
	)
	ctx.EventManager().EmitEvent(event)

	if k.hooks != nil {
		k.hooks.AfterLeveragelpPositionModified(ctx, ammPool, pool)
	}

	return &types.MsgOpenResponse{}, nil
}

func (k Keeper) ProcessOpenLong(ctx sdk.Context, position *types.Position, leverage sdk.Dec, collateralAmountDec sdk.Dec, poolId uint64, msg *types.MsgOpen) (*types.Position, error) {
	// Fetch the pool associated with the given pool ID.
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	// Check if the pool is enabled.
	if !k.IsPoolEnabled(ctx, poolId) {
		return nil, sdkerrors.Wrap(types.ErrPositionDisabled, fmt.Sprintf("poolId: %d", poolId))
	}

	// Fetch the corresponding AMM (Automated Market Maker) pool.
	ammPool, err := k.GetAmmPool(ctx, poolId)
	if err != nil {
		return nil, err
	}

	if msg.CollateralAsset != ptypes.BaseCurrency {
		return nil, types.ErrOnlyBaseCurrencyAllowed
	}

	// Calculate the leveraged amount based on the collateral provided and the leverage.
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(leverage).TruncateInt().Int64())

	// send collateral coins to Position address from Position owner address
	positionOwner := sdk.MustAccAddressFromBech32(position.Address)
	err = k.bankKeeper.SendCoins(ctx, positionOwner, position.GetPositionAddress(), sdk.Coins{sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)})
	if err != nil {
		return nil, err
	}
	leverageCoin := sdk.NewCoin(msg.CollateralAsset, leveragedAmount)

	// borrow leveragedAmount - collateralAmount
	borrowCoin := sdk.NewCoin(msg.CollateralAsset, leveragedAmount.Sub(msg.CollateralAmount))
	err = k.stableKeeper.Borrow(ctx, position.GetPositionAddress(), borrowCoin)
	if err != nil {
		return nil, err
	}

	_, shares, err := k.amm.JoinPoolNoSwap(ctx, position.GetPositionAddress(), poolId, sdk.OneInt(), sdk.Coins{leverageCoin}, true)
	if err != nil {
		return nil, err
	}

	// Update the pool health.
	pool.LeveragedLpAmount = pool.LeveragedLpAmount.Add(shares)
	if err = k.UpdatePoolHealth(ctx, &pool); err != nil {
		return nil, err
	}

	// Update the Position health.
	lr, err := k.GetPositionHealth(ctx, *position, ammPool)
	if err != nil {
		return nil, err
	}

	// Check if the Position is unhealthy
	safetyFactor := k.GetSafetyFactor(ctx)
	if lr.LTE(safetyFactor) {
		return nil, types.ErrPositionUnhealthy
	}

	// Set Position
	position.LeveragedLpAmount = position.LeveragedLpAmount.Add(shares)
	position.Liabilities = position.Liabilities.Add(borrowCoin.Amount)
	position.PositionHealth = lr
	k.SetPosition(ctx, position)

	return position, nil
}
