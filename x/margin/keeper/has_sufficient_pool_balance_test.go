package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/stretchr/testify/assert"
)

func TestHasSufficientPoolBalance_SufficientBalance(t *testing.T) {
	// Setup the keeper
	k := keeper.Keeper{}

	ctx := sdk.Context{} // mock or setup a context

	// Define the ammPool with assets
	ammPool := types.Pool{
		PoolAssets: []types.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: sdk.NewInt(100),
				},
			},
		},
	}

	borrowAsset := "testAsset"
	requiredAmount := sdk.NewInt(50)

	// Run the function
	hasBalance := k.HasSufficientPoolBalance(ctx, ammPool, borrowAsset, requiredAmount)

	// Expect that there is a sufficient balance
	assert.True(t, hasBalance)
}

func TestHasSufficientPoolBalance_InsufficientBalance(t *testing.T) {
	// Setup the keeper
	k := keeper.Keeper{}

	ctx := sdk.Context{} // mock or setup a context

	// Define the ammPool with assets
	ammPool := types.Pool{
		PoolAssets: []types.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: sdk.NewInt(100),
				},
			},
		},
	}

	borrowAsset := "testAsset"
	requiredAmount := sdk.NewInt(150)

	// Run the function
	hasBalance := k.HasSufficientPoolBalance(ctx, ammPool, borrowAsset, requiredAmount)

	// Expect that there is an insufficient balance
	assert.False(t, hasBalance)
}
