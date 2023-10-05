package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtype "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/keeper"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/stretchr/testify/assert"
)

func TestGetAmmPoolBalance_GetAmmPoolBalanceAvailable(t *testing.T) {
	// Setup the keeper
	k := keeper.Keeper{}

	ctx := sdk.Context{} // mock or setup a context

	// Define the ammPool with assets
	ammPool := ammtype.Pool{
		PoolAssets: []ammtype.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: sdk.NewInt(100),
				},
			},
		},
	}

	borrowAsset := "testAsset"

	// Run the function
	balance, err := k.GetAmmPoolBalance(ctx, ammPool, borrowAsset)

	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is 100 balance
	assert.Equal(t, balance, sdk.NewInt(100))
}

func TestGetAmmPoolBalance_GetAmmPoolBalanceUnavailable(t *testing.T) {
	// Setup the keeper
	k := keeper.Keeper{}

	ctx := sdk.Context{} // mock or setup a context

	// Define the ammPool with assets
	ammPool := ammtype.Pool{
		PoolAssets: []ammtype.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: sdk.NewInt(100),
				},
			},
		},
	}

	borrowAsset := "testAsset2"

	// Run the function
	balance, err := k.GetAmmPoolBalance(ctx, ammPool, borrowAsset)

	// Expect that there is an insufficient balance
	assert.True(t, errors.Is(err, sdkerrors.Wrap(types.ErrBalanceNotAvailable, "Balance not available")))
	assert.Equal(t, balance, sdk.ZeroInt())
}
