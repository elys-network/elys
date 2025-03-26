package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeEvtPoolJoined      = "pool_joined"
	TypeEvtPoolExited      = "pool_exited"
	TypeEvtPoolCreated     = "pool_created"
	TypeEvtTokenSwapped    = "token_swapped"
	TypeEvtTokenSwappedFee = "token_swapped_fee"

	AttributeValueCategory = ModuleName
	AttributeKeyPoolId     = "pool_id"
	AttributeKeyTokensIn   = "tokens_in"
	AttributeKeyTokensOut  = "tokens_out"
	AttributeKeyRecipient  = "recipient"

	AttributeKeySwapFee           = "swap_fee"
	AttributeKeySlippage          = "slippage"
	AttributeKeyWeightRecoveryFee = "weight_recovery_fee"
	AttributeKeyProvidedBonusFee  = "provided_bonus_fee"
	AttributeTakerFees            = "taker_fees"
)

func EmitSwapEvent(ctx sdk.Context, sender, recipient sdk.AccAddress, poolId uint64, input sdk.Coins, output sdk.Coins) {
	ctx.EventManager().EmitEvents(sdk.Events{
		NewSwapEvent(sender, recipient, poolId, input, output),
	})
}

func EmitSwapFeesCollectedEvent(ctx sdk.Context, swapFee string, slippage string, weightRecoveryFee string, providedBonusFee string, takerFees string) {
	ctx.EventManager().EmitEvents(sdk.Events{
		NewSwapFeeEvent(swapFee, slippage, weightRecoveryFee, providedBonusFee, takerFees),
	})
}

func EmitAddLiquidityEvent(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, liquidity sdk.Coins) {
	ctx.EventManager().EmitEvents(sdk.Events{
		NewAddLiquidityEvent(sender, poolId, liquidity),
	})
}

func EmitRemoveLiquidityEvent(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, liquidity sdk.Coins) {
	ctx.EventManager().EmitEvents(sdk.Events{
		NewRemoveLiquidityEvent(sender, poolId, liquidity),
	})
}

func NewSwapEvent(sender, recipient sdk.AccAddress, poolId uint64, input sdk.Coins, output sdk.Coins) sdk.Event {
	return sdk.NewEvent(
		TypeEvtTokenSwapped,
		sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(AttributeKeyRecipient, recipient.String()),
		sdk.NewAttribute(AttributeKeyPoolId, strconv.FormatUint(poolId, 10)),
		sdk.NewAttribute(AttributeKeyTokensIn, input.String()),
		sdk.NewAttribute(AttributeKeyTokensOut, output.String()),
	)
}

func NewSwapFeeEvent(swapFee string, slippage string, weightRecoveryFee string, providedBonusFee string, takerFees string) sdk.Event {
	return sdk.NewEvent(
		TypeEvtTokenSwappedFee,
		sdk.NewAttribute("denom", "USD"),
		sdk.NewAttribute(AttributeKeySwapFee, swapFee),
		sdk.NewAttribute(AttributeKeySlippage, slippage),
		sdk.NewAttribute(AttributeKeyWeightRecoveryFee, weightRecoveryFee),
		sdk.NewAttribute(AttributeKeyProvidedBonusFee, providedBonusFee),
		sdk.NewAttribute(AttributeTakerFees, takerFees),
	)
}

func NewAddLiquidityEvent(sender sdk.AccAddress, poolId uint64, liquidity sdk.Coins) sdk.Event {
	return sdk.NewEvent(
		TypeEvtPoolJoined,
		sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(AttributeKeyPoolId, strconv.FormatUint(poolId, 10)),
		sdk.NewAttribute(AttributeKeyTokensIn, liquidity.String()),
	)
}

func NewRemoveLiquidityEvent(sender sdk.AccAddress, poolId uint64, liquidity sdk.Coins) sdk.Event {
	return sdk.NewEvent(
		TypeEvtPoolExited,
		sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(AttributeKeyPoolId, strconv.FormatUint(poolId, 10)),
		sdk.NewAttribute(AttributeKeyTokensOut, liquidity.String()),
	)
}
