package keeper

import (
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) OpenLong(ctx sdk.Context, msg *types.MsgOpen) (*types.Position, error) {
	// Initialize a new Leveragelp Trading Position (Position).
	if msg.Leverage.LT(sdk.OneDec()) {
		return nil, types.ErrLeverageTooSmall
	}
	position := types.NewPosition(msg.Creator, sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount), msg.AmmPoolId)
	position.Id = k.GetPositionCount(ctx) + 1
	position.StopLossPrice = msg.StopLossPrice
	k.SetPositionCount(ctx, position.Id)

	openCount := k.GetOpenPositionCount(ctx)
	k.SetOpenPositionCount(ctx, openCount+1)

	// Call the function to process the open long logic.
	return k.ProcessOpenLong(ctx, position, msg.AmmPoolId, msg)
}

func (k Keeper) OpenConsolidate(ctx sdk.Context, position *types.Position, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	if msg.Leverage.LT(sdk.OneDec()) {
		return nil, types.ErrLeverageTooSmall
	}

	poolId := position.AmmPoolId

	position.Collateral = position.Collateral.Add(sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount))

	position, err := k.ProcessOpenLong(ctx, position, poolId, msg)
	if err != nil {
		return nil, err
	}

	if k.hooks != nil {
		err := k.hooks.AfterLeverageLpPositionOpenConsolidate(ctx, sdk.MustAccAddressFromBech32(msg.Creator))
		if err != nil {
			return nil, err
		}

	}

	event := sdk.NewEvent(types.EventOpen,
		sdk.NewAttribute("id", strconv.FormatInt(int64(position.Id), 10)),
		sdk.NewAttribute("address", position.Address),
		sdk.NewAttribute("collateral", position.Collateral.String()),
		sdk.NewAttribute("liabilities", position.Liabilities.String()),
		sdk.NewAttribute("health", position.PositionHealth.String()),
	)
	ctx.EventManager().EmitEvent(event)

	return &types.MsgOpenResponse{}, nil
}

func (k Keeper) ProcessOpenLong(ctx sdk.Context, position *types.Position, poolId uint64, msg *types.MsgOpen) (*types.Position, error) {
	collateralAmountDec := sdk.NewDecFromInt(msg.CollateralAmount)
	// Determine the maximum leverage available and compute the effective leverage to be used.
	maxLeverage := k.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)

	// Fetch the pool associated with the given pool ID.
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	// Check if the pool is enabled.
	if !pool.Enabled {
		return nil, errorsmod.Wrap(types.ErrPositionDisabled, fmt.Sprintf("poolId: %d", poolId))
	}

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	if msg.CollateralAsset != baseCurrency {
		return nil, types.ErrOnlyBaseCurrencyAllowed
	}

	// Calculate the leveraged amount based on the collateral provided and the leverage.
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(leverage).TruncateInt().Int64())

	// send collateral coins to Position address from Position owner address
	positionOwner := sdk.MustAccAddressFromBech32(position.Address)
	err := k.bankKeeper.SendCoins(ctx, positionOwner, position.GetPositionAddress(), sdk.Coins{sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)})
	if err != nil {
		return nil, err
	}
	leverageCoin := sdk.NewCoin(msg.CollateralAsset, leveragedAmount)

	// borrow leveragedAmount - collateralAmount
	borrowCoin := sdk.NewCoin(msg.CollateralAsset, leveragedAmount.Sub(msg.CollateralAmount))
	if borrowCoin.Amount.IsPositive() {
		err = k.stableKeeper.Borrow(ctx, position.GetPositionAddress(), borrowCoin)
		if err != nil {
			return nil, err
		}
	}

	_, shares, err := k.amm.JoinPoolNoSwap(ctx, position.GetPositionAddress(), poolId, sdk.OneInt(), sdk.Coins{leverageCoin})
	if err != nil {
		return nil, err
	}

	// Update the pool health.
	pool.LeveragedLpAmount = pool.LeveragedLpAmount.Add(shares)
	k.UpdatePoolHealth(ctx, &pool)

	// Get the Position health.
	lr, err := k.GetPositionHealth(ctx, *position)
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
	position.StopLossPrice = msg.StopLossPrice

	k.SetPosition(ctx, position)

	return position, nil
}
