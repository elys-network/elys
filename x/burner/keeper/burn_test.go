package keeper_test

import (
	"cosmossdk.io/math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/burner/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestShouldBurnTokens(t *testing.T) {
	k, ctx, _ := keepertest.BurnerKeeper(t)

	params := types.Params{
		EpochIdentifier: "test_epoch",
	}
	k.SetParams(ctx, &params)

	require.True(t, k.ShouldBurnTokens(ctx, "test_epoch"))
	require.False(t, k.ShouldBurnTokens(ctx, "invalid_epoch"))
}

func TestBurnTokensForAllDenoms(t *testing.T) {
	k, ctx, bankKeeper := keepertest.BurnerKeeper(t)
	// addr := sdk.AccAddress([]byte("module"))

	// Set up some balances
	balances := []struct {
		denom  string
		amount int64
	}{
		{"denom1", 100},
		{"denom2", 200},
		{"denom3", 0}, // zero balance should be ignored
	}

	bankKeeper.EXPECT().IterateAllDenomMetaData(ctx, mock.Anything).Run(func(ctx sdk.Context, callback func(metadata banktypes.Metadata) bool) {
		callback(banktypes.Metadata{Base: balances[0].denom})
		callback(banktypes.Metadata{Base: balances[1].denom})
	}).Once()

	bankKeeper.EXPECT().GetBalance(ctx, types.GetZeroAddress(), balances[0].denom).Return(sdk.NewCoin(balances[0].denom, math.NewInt(balances[0].amount))).Once()
	bankKeeper.EXPECT().GetBalance(ctx, types.GetZeroAddress(), balances[1].denom).Return(sdk.NewCoin(balances[1].denom, math.NewInt(balances[1].amount))).Once()

	bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, types.GetZeroAddress(), types.ModuleName, sdk.NewCoins(sdk.NewCoin(balances[0].denom, math.NewInt(balances[0].amount)))).Return(nil).Once()
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, types.GetZeroAddress(), types.ModuleName, sdk.NewCoins(sdk.NewCoin(balances[1].denom, math.NewInt(balances[1].amount)))).Return(nil).Once()

	bankKeeper.EXPECT().BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(balances[0].denom, math.NewInt(balances[0].amount)))).Return(nil)
	bankKeeper.EXPECT().BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(balances[1].denom, math.NewInt(balances[1].amount)))).Return(nil)

	// Burn the tokens
	err := k.BurnTokensForAllDenoms(ctx)
	require.NoError(t, err)
}
