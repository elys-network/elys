package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeEvtPoolJoined   = "pool_joined"
	TypeEvtPoolCreated  = "pool_created"
	TypeEvtTokenSwapped = "token_swapped"

	AttributeValueCategory = ModuleName
	AttributeKeyPoolId     = "pool_id"
	AttributeKeyTokensIn   = "tokens_in"
	AttributeKeyTokensOut  = "tokens_out"
)

func EmitSwapEvent(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, input sdk.Coins, output sdk.Coins) {
	ctx.EventManager().EmitEvents(sdk.Events{
		newSwapEvent(sender, poolId, input, output),
	})
}

func EmitAddLiquidityEvent(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, liquidity sdk.Coins) {
	ctx.EventManager().EmitEvents(sdk.Events{
		newAddLiquidityEvent(sender, poolId, liquidity),
	})
}

func newSwapEvent(sender sdk.AccAddress, poolId uint64, input sdk.Coins, output sdk.Coins) sdk.Event {
	return sdk.NewEvent(
		TypeEvtTokenSwapped,
		sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(AttributeKeyPoolId, strconv.FormatUint(poolId, 10)),
		sdk.NewAttribute(AttributeKeyTokensIn, input.String()),
		sdk.NewAttribute(AttributeKeyTokensOut, output.String()),
	)
}

func newAddLiquidityEvent(sender sdk.AccAddress, poolId uint64, liquidity sdk.Coins) sdk.Event {
	return sdk.NewEvent(
		TypeEvtPoolJoined,
		sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(AttributeKeyPoolId, strconv.FormatUint(poolId, 10)),
		sdk.NewAttribute(AttributeKeyTokensIn, liquidity.String()),
	)
}
