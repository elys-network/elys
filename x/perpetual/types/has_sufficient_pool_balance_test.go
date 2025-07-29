package types_test

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/stretchr/testify/assert"
)

func TestHasSufficientPoolBalance_SufficientBalance(t *testing.T) {
	// Define the ammPool with assets
	ammPool := ammtypes.Pool{
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: math.NewInt(100),
				},
			},
		},
	}

	borrowAsset := "testAsset"
	requiredAmount := math.NewInt(50)

	// Run the function
	hasBalance := types.HasSufficientPoolBalance(ammPool, borrowAsset, requiredAmount)

	// Expect that there is a sufficient balance
	assert.True(t, hasBalance)
}

func TestHasSufficientPoolBalance_InsufficientBalance(t *testing.T) {
	// Define the ammPool with assets
	ammPool := ammtypes.Pool{
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: math.NewInt(100),
				},
			},
		},
	}

	borrowAsset := "testAsset"
	requiredAmount := math.NewInt(150)

	// Run the function
	hasBalance := types.HasSufficientPoolBalance(ammPool, borrowAsset, requiredAmount)

	// Expect that there is an insufficient balance
	assert.False(t, hasBalance)
}
