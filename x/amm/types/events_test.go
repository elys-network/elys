package types_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestEvents(t *testing.T) {
	sender := sdk.AccAddress([]byte("sender"))
	poolID := uint64(123)
	input := sdk.NewCoins(sdk.NewInt64Coin("token1", 100))
	output := sdk.NewCoins(sdk.NewInt64Coin("token2", 50))

	swapEvent := types.NewSwapEvent(sender, sender, poolID, input, output)
	expectedSwapEvent := sdk.NewEvent(
		types.TypeEvtTokenSwapped,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(types.AttributeKeyRecipient, sender.String()),
		sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(poolID, 10)),
		sdk.NewAttribute(types.AttributeKeyTokensIn, input.String()),
		sdk.NewAttribute(types.AttributeKeyTokensOut, output.String()),
	)
	require.Equal(t, expectedSwapEvent, swapEvent)

	liquidity := sdk.NewCoins(sdk.NewInt64Coin("token3", 200))

	addLiquidityEvent := types.NewAddLiquidityEvent(sender, poolID, liquidity)
	expectedAddLiquidityEvent := sdk.NewEvent(
		types.TypeEvtPoolJoined,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(poolID, 10)),
		sdk.NewAttribute(types.AttributeKeyTokensIn, liquidity.String()),
	)
	require.Equal(t, expectedAddLiquidityEvent, addLiquidityEvent)

	removeLiquidityEvent := types.NewRemoveLiquidityEvent(sender, poolID, liquidity)
	expectedRemoveLiquidityEvent := sdk.NewEvent(
		types.TypeEvtPoolExited,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(poolID, 10)),
		sdk.NewAttribute(types.AttributeKeyTokensOut, liquidity.String()),
	)
	require.Equal(t, expectedRemoveLiquidityEvent, removeLiquidityEvent)
}
