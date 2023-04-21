package keeper_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/burner/types"
	"github.com/stretchr/testify/require"
)

func TestShouldBurnTokens(t *testing.T) {
	k, ctx := keepertest.BurnerKeeper(t)

	params := types.Params{
		EpochIdentifier: "test_epoch",
	}
	k.SetParams(ctx, params)

	require.True(t, k.ShouldBurnTokens(ctx, "test_epoch"))
	require.False(t, k.ShouldBurnTokens(ctx, "invalid_epoch"))
}

// func TestBurnTokensForAllDenoms(t *testing.T) {
// 	k, ctx := keepertest.BurnerKeeper(t)
// 	addr := sdk.AccAddress([]byte("module"))

// 	// Set up some balances
// 	balances := []struct {
// 		denom  string
// 		amount int64
// 	}{
// 		{"denom1", 100},
// 		{"denom2", 200},
// 		{"denom3", 0}, // zero balance should be ignored
// 	}

// 	for _, b := range balances {
// 		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(b.denom, sdk.NewInt(b.amount))))
// 		require.NoError(t, err)

// 		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(sdk.NewCoin(b.denom, sdk.NewInt(b.amount))))
// 		require.NoError(t, err)
// 	}

// 	// Burn the tokens
// 	err := k.BurnTokensForAllDenoms(ctx)
// 	require.NoError(t, err)

// 	// Verify that the tokens have been burned
// 	for _, b := range balances {
// 		balance := k.bankKeeper.GetBalance(ctx, addr, b.denom)
// 		require.Equal(t, sdk.ZeroInt(), balance.Amount)
// 	}
// }
