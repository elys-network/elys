package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/stretchr/testify/require"
)

func TestCalcSwapEstimationByDenom(t *testing.T) {
	k, ctx, accountedPoolKeeper, _ := keepertest.AmmKeeper(t)
	SetupMockPools(k, ctx)

	// Test with amount denom equal to denomIn
	amount := sdk.NewCoin("denom1", sdk.NewInt(100))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(1), "denom1").Return(sdk.NewInt(1000))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(1), "denom2").Return(sdk.NewInt(1000))
	inRoute, outRoute, tokenOut, spotPrice, _, _, _, err := k.CalcSwapEstimationByDenom(ctx, amount, "denom1", "denom2", "baseCurrency", sdk.ZeroDec())
	require.NoError(t, err)
	require.NotNil(t, inRoute)
	require.Nil(t, outRoute)
	require.NotZero(t, tokenOut)
	require.NotZero(t, spotPrice)

	// Test with amount denom equal to denomOut
	amount = sdk.NewCoin("denom2", sdk.NewInt(100))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(1), "denom2").Return(sdk.NewInt(1000))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(1), "denom1").Return(sdk.NewInt(1000))
	inRoute, outRoute, tokenOut, spotPrice, _, _, _, err = k.CalcSwapEstimationByDenom(ctx, amount, "denom1", "denom2", "baseCurrency", sdk.ZeroDec())
	require.NoError(t, err)
	require.Nil(t, inRoute)
	require.NotNil(t, outRoute)
	require.NotZero(t, tokenOut)
	require.NotZero(t, spotPrice)

	// Test with invalid amount denom
	amount = sdk.NewCoin("invalid", sdk.NewInt(1000))
	_, _, _, _, _, _, _, err = k.CalcSwapEstimationByDenom(ctx, amount, "denom1", "denom2", "baseCurrency", sdk.ZeroDec())
	require.Error(t, err)
}
