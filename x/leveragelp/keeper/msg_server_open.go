package keeper

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	stabletypes "github.com/elys-network/elys/v6/x/stablestake/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k msgServer) Open(goCtx context.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	enabledPools := k.GetParams(ctx).EnabledPools
	found := slices.Contains(enabledPools, msg.AmmPoolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolNotEnabled, fmt.Sprintf("poolId: %d", msg.AmmPoolId))
	}

	return k.Keeper.Open(ctx, msg)
}

func (k Keeper) Open(ctx sdk.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	if err := k.CheckUserAuthorization(ctx, msg); err != nil {
		return nil, err
	}
	moduleAddr := authtypes.NewModuleAddress(stabletypes.ModuleName)

	borrowPool, found := k.stableKeeper.GetPoolByDenom(ctx, msg.CollateralAsset)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolNotCreatedForBorrow, fmt.Sprintf("Asset: %s", msg.CollateralAsset))
	}

	depositDenom := borrowPool.GetDepositDenom()
	balance := k.bankKeeper.GetBalance(ctx, moduleAddr, depositDenom)
	borrowed := osmomath.BigDecFromSDKInt(borrowPool.NetAmount.Sub(balance.Amount))
	borrowRatio := osmomath.ZeroBigDec()
	if borrowPool.NetAmount.GT(sdkmath.ZeroInt()) {
		borrowRatio = borrowed.Add(osmomath.BigDecFromDec(msg.Leverage).Mul(osmomath.BigDecFromSDKInt(msg.CollateralAmount))).
			Quo(borrowPool.GetBigDecNetAmount())
	}

	ammPool, found := k.amm.GetPool(ctx, msg.AmmPoolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", msg.AmmPoolId))
	}

	found = false
	for _, asset := range ammPool.PoolAssets {
		if asset.Token.Denom == msg.CollateralAsset {
			found = true
			break
		}
	}
	if !found {
		return nil, errorsmod.Wrap(types.ErrAssetNotSupported, fmt.Sprintf("Asset: %s", msg.CollateralAsset))
	}

	if borrowRatio.GTE(borrowPool.GetBigDecMaxLeverageRatio()) {
		return nil, errors.New("stable stake pool max borrow capacity used up")
	}

	// Check if it is the same direction position for the same trader.
	if position := k.CheckSamePosition(ctx, msg); position != nil {
		response, err := k.OpenConsolidate(ctx, position, msg)
		if err != nil {
			return nil, err
		}
		if err = k.CheckMaxLeverageRatio(ctx, msg.AmmPoolId); err != nil {
			return nil, err
		}
		return response, nil
	}

	if err := k.CheckMaxOpenPositions(ctx, msg.AmmPoolId); err != nil {
		return nil, err
	}

	position, err := k.OpenLong(ctx, msg, borrowPool.Id)
	if err != nil {
		return nil, err
	}

	if err = k.CheckMaxLeverageRatio(ctx, msg.AmmPoolId); err != nil {
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
		sdk.NewAttribute("poolId", strconv.FormatUint(position.AmmPoolId, 10)),
		sdk.NewAttribute("collateral", position.Collateral.String()),
		sdk.NewAttribute("liabilities", position.Liabilities.String()),
		sdk.NewAttribute("health", position.PositionHealth.String()),
		sdk.NewAttribute("leverage_lp_amount", position.LeveragedLpAmount.String()),
	)
	ctx.EventManager().EmitEvent(event)

	return &types.MsgOpenResponse{}, nil
}
