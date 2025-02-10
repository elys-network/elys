package keeper

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k msgServer) Open(goCtx context.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	enabledPools := k.GetParams(ctx).EnabledPools
	found := false
	for _, poolId := range enabledPools {
		if poolId == msg.AmmPoolId {
			found = true
			break
		}
	}
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolNotEnabled, fmt.Sprintf("poolId: %d", msg.AmmPoolId))
	}

	return k.Keeper.Open(ctx, msg)
}

func (k Keeper) Open(ctx sdk.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	if err := k.CheckUserAuthorization(ctx, msg); err != nil {
		return nil, err
	}
	stableStakeParams := k.stableKeeper.GetParams(ctx)
	moduleAddr := authtypes.NewModuleAddress(stabletypes.ModuleName)

	depositDenom := k.stableKeeper.GetDepositDenom(ctx)

	balance := k.bankKeeper.GetBalance(ctx, moduleAddr, depositDenom)
	borrowed := stableStakeParams.TotalValue.Sub(balance.Amount)
	borrowRatio := sdkmath.LegacyZeroDec()
	if stableStakeParams.TotalValue.GT(sdkmath.ZeroInt()) {
		borrowRatio = borrowed.ToLegacyDec().Add(msg.Leverage.Mul(msg.CollateralAmount.ToLegacyDec())).
			Quo(stableStakeParams.TotalValue.ToLegacyDec())
	}

	if borrowRatio.GTE(stableStakeParams.MaxLeverageRatio) {
		return nil, errors.New("stable stake pool max borrow capacity used up")
	}

	// Check if it is the same direction position for the same trader.
	if position, err := k.CheckSamePosition(ctx, msg); position != nil {
		if err != nil {
			return nil, err
		}
		response, err := k.OpenConsolidate(ctx, position, msg)
		if err != nil {
			return nil, err
		}
		if err = k.CheckPoolHealth(ctx, msg.AmmPoolId); err != nil {
			return nil, err
		}
		return response, nil
	}

	if err := k.CheckMaxOpenPositions(ctx); err != nil {
		return nil, err
	}

	position, err := k.OpenLong(ctx, msg)
	if err != nil {
		return nil, err
	}

	if err = k.CheckPoolHealth(ctx, msg.AmmPoolId); err != nil {
		return nil, err
	}

	if k.hooks != nil {
		// ammPool will have updated values for opening position
		ammPool, found := k.amm.GetPool(ctx, msg.AmmPoolId)
		if !found {
			return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", msg.AmmPoolId))
		}
		err = k.hooks.AfterLeverageLpPositionOpen(ctx, sdk.MustAccAddressFromBech32(msg.Creator), ammPool)
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
		sdk.NewAttribute("leverage_lp_amount", position.LeveragedLpAmount.String()),
	)
	ctx.EventManager().EmitEvent(event)

	return &types.MsgOpenResponse{}, nil
}
