package keeper

import (
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/amm/types"
)

// EndBlocker of amm module
func (k Keeper) EndBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	swapInRequests := k.GetAllSwapExactAmountInRequests(ctx)
	for _, msg := range swapInRequests {
		sender, err := sdk.AccAddressFromBech32(msg.Sender)
		if err != nil {
			continue
		}

		cacheCtx, write := ctx.CacheContext()
		_, err = k.RouteExactAmountIn(cacheCtx, sender, msg.Routes, msg.TokenIn, math.Int(msg.TokenOutMinAmount))
		if err != nil {
			continue
		}
		write()

		// Swap event is handled elsewhere
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			),
		})

	}
	swapOutRequests := k.GetAllSwapExactAmountOutRequests(ctx)
	for _, msg := range swapOutRequests {
		sender, err := sdk.AccAddressFromBech32(msg.Sender)
		if err != nil {
			continue
		}

		cacheCtx, write := ctx.CacheContext()
		_, err = k.RouteExactAmountOut(cacheCtx, sender, msg.Routes, msg.TokenInMaxAmount, msg.TokenOut)
		if err != nil {
			continue
		}
		write()

		// Swap event is handled elsewhere
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			),
		})
	}
}
