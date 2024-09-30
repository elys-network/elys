package keeper

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k msgServer) ClosePositions(goCtx context.Context, msg *types.MsgClosePositions) (*types.MsgClosePositionsResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	baseCurrency, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, nil
	}

	// Handle toPay
	err := k.HandleToPay(ctx)
	if err != nil {
		ctx.Logger().Error(err.Error())
	}

	// Handle liquidations
	liqLog := []string{}
	for _, val := range msg.Liquidate {
		owner := sdk.MustAccAddressFromBech32(val.Address)
		position, err := k.GetMTP(ctx, owner, val.Id)
		if err != nil {
			continue
		}

		pool, poolFound := k.GetPool(ctx, position.AmmPoolId)
		if !poolFound {
			continue
		}
		ammPool, poolErr := k.GetAmmPool(ctx, position.AmmPoolId, position.TradingAsset)
		if poolErr != nil {
			continue
		}

		err = k.CheckAndLiquidateUnhealthyPosition(ctx, &position, pool, ammPool, baseCurrency.Denom, baseCurrency.Decimals)
		if err != nil {
			// Add log about error or not liquidated
			liqLog = append(liqLog, fmt.Sprintf("Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		}
	}

	//Handle StopLoss
	closeLog := []string{}
	for _, val := range msg.StopLoss {
		owner := sdk.MustAccAddressFromBech32(val.Address)
		position, err := k.GetMTP(ctx, owner, val.Id)
		if err != nil {
			continue
		}

		pool, poolFound := k.GetPool(ctx, position.AmmPoolId)
		if !poolFound {
			continue
		}
		ammPool, poolErr := k.GetAmmPool(ctx, position.AmmPoolId, position.TradingAsset)
		if poolErr != nil {
			continue
		}

		err = k.CheckAndCloseAtStopLoss(ctx, &position, pool, ammPool, baseCurrency.Denom, baseCurrency.Decimals)
		if err != nil {
			// Add log about error or not closed
			closeLog = append(closeLog, fmt.Sprintf("Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		}
	}

	// TODO: Handle stop loss
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClosePositions,
		sdk.NewAttribute("liquidations", strings.Join(liqLog, "\n")),
		sdk.NewAttribute("stop_loss", strings.Join(closeLog, "\n")),
	))

	return &types.MsgClosePositionsResponse{}, nil
}
