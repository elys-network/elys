package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"strconv"
)

func (k Keeper) EmitCloseEvent(ctx sdk.Context, trigger string, position types.Position, closingRatio math.LegacyDec, totalLpAmountToClose math.Int, coinsForAmm sdk.Coins, repayAmount math.Int, userReturnTokens sdk.Coins, exitFeeOnClosingPosition osmomath.BigDec, stopLossReached bool, exitSlippageFee, swapFee, takerFee osmomath.BigDec) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
		sdk.NewAttribute("id", strconv.FormatUint(position.Id, 10)),
		sdk.NewAttribute("address", position.Address),
		sdk.NewAttribute("poolId", strconv.FormatUint(position.AmmPoolId, 10)),
		sdk.NewAttribute("closing_ratio", closingRatio.String()),
		sdk.NewAttribute("lp_amount_closed", totalLpAmountToClose.String()),
		sdk.NewAttribute("coins_to_amm", coinsForAmm.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("user_return_tokens", userReturnTokens.String()),
		sdk.NewAttribute("exit_fee", exitFeeOnClosingPosition.String()),
		sdk.NewAttribute("health", position.PositionHealth.String()),
		sdk.NewAttribute("stop_loss_reached", strconv.FormatBool(stopLossReached)),
		sdk.NewAttribute("exit_slippage_fee", exitSlippageFee.String()),
		sdk.NewAttribute("exit_swap_fee", swapFee.String()),
		sdk.NewAttribute("exit_taker_fee", takerFee.String()),
		sdk.NewAttribute("trigger", trigger),
	))
}
