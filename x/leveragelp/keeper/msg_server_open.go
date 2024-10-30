package keeper

import (
	"context"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k msgServer) Open(goCtx context.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return k.Keeper.Open(ctx, msg)
}

func (k Keeper) Open(ctx sdk.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	if err := k.CheckUserAuthorization(ctx, msg); err != nil {
		return nil, err
	}
	params := k.stableKeeper.GetParams(ctx)
	moduleAddr := authtypes.NewModuleAddress(stabletypes.ModuleName)

	depositDenom := k.stableKeeper.GetDepositDenom(ctx)

	balance := k.bankKeeper.GetBalance(ctx, moduleAddr, depositDenom)
	borrowed := params.TotalValue.Sub(balance.Amount)
	borrowRatio := sdk.ZeroDec()
	if params.TotalValue.GT(sdk.ZeroInt()) {
		borrowRatio = borrowed.ToLegacyDec().Add(msg.Leverage.Mul(msg.CollateralAmount.ToLegacyDec())).
			Quo(params.TotalValue.ToLegacyDec())
	}

	var poolLeveragelpRatio sdk.Dec
	pool, found := k.GetPool(ctx, msg.AmmPoolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", msg.AmmPoolId))
	}
	amm_pool, found := k.amm.GetPool(ctx, msg.AmmPoolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", msg.AmmPoolId))
	}
	poolLeveragelpRatio = pool.LeveragedLpAmount.ToLegacyDec().Quo(amm_pool.TotalShares.Amount.ToLegacyDec())

	if poolLeveragelpRatio.GTE(pool.MaxLeveragelpRatio) || borrowRatio.GTE(params.MaxLeverageRatio) {
		return nil, errorsmod.Wrap(types.ErrMaxLeverageLpExists, "no new position can be open")
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
		err := k.hooks.AfterLeverageLpPositionOpen(ctx, sdk.MustAccAddressFromBech32(msg.Creator))
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
