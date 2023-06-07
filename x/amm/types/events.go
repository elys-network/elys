package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
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
