package types_test

import (
	"errors"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/assert"
)

func TestPool_UpdateLiabilitiesValid(t *testing.T) {
	// ctx := sdk.Context{} // mock or setup a context

	// Define the perpetual pool with assets
	ammPool := ammtypes.Pool{
		PoolId: 1,
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: sdk.NewInt(100),
				},
			},
		},
	}
	pool := types.NewPool(ammPool)
	pool.PoolAssetsLong = []types.PoolAsset{
		{
			Liabilities: sdk.NewInt(0),
			Custody:     sdk.NewInt(0),
			AssetDenom:  "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 150.
	denom := "testAsset"
	err := pool.UpdateLiabilities(denom, sdk.NewInt(100), true, types.Position_LONG)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is 100 liabilities
	assert.Equal(t, pool.PoolAssetsLong[0].Liabilities, sdk.NewInt(100))
	err = pool.UpdateLiabilities(denom, sdk.NewInt(150), false, types.Position_LONG)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is -50 liabilities
	assert.Equal(t, pool.PoolAssetsLong[0].Liabilities, sdk.NewInt(-50))
}

func TestPool_UpdateLiabilitiesInvalid(t *testing.T) {
	// ctx := sdk.Context{} // mock or setup a context

	// Define the perpetual pool with assets
	ammPool := ammtypes.Pool{
		PoolId: 1,
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: sdk.NewInt(100),
				},
			},
		},
	}
	pool := types.NewPool(ammPool)
	pool.PoolAssetsLong = []types.PoolAsset{
		{
			Liabilities: sdk.NewInt(0),
			Custody:     sdk.NewInt(0),
			AssetDenom:  "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 50.
	denom := "testAsset2"
	err := pool.UpdateLiabilities(denom, sdk.NewInt(100), true, types.Position_LONG)
	// Expect that there is invalid asset denom error.
	assert.True(t, errors.Is(err, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")))

	// Expect that there is still 0 liabilities
	assert.Equal(t, pool.PoolAssetsLong[0].Liabilities, sdk.NewInt(0))
}

func TestPool_UpdateTakeProfitLiabilitiesValid(t *testing.T) {
	// ctx := sdk.Context{} // mock or setup a context

	// Define the perpetual pool with assets
	ammPool := ammtypes.Pool{
		PoolId: 1,
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: sdk.NewInt(100),
				},
			},
		},
	}
	pool := types.NewPool(ammPool)
	pool.PoolAssetsLong = []types.PoolAsset{
		{
			Liabilities:           sdk.NewInt(0),
			TakeProfitLiabilities: sdk.NewInt(0),
			Custody:               sdk.NewInt(0),
			AssetDenom:            "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 150.
	denom := "testAsset"
	err := pool.UpdateTakeProfitLiabilities(denom, sdk.NewInt(100), true, types.Position_LONG)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is 100 liabilities
	assert.Equal(t, pool.PoolAssetsLong[0].TakeProfitLiabilities, sdk.NewInt(100))
	err = pool.UpdateTakeProfitLiabilities(denom, sdk.NewInt(150), false, types.Position_LONG)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is -50 liabilities
	assert.Equal(t, pool.PoolAssetsLong[0].TakeProfitLiabilities, sdk.NewInt(-50))
}

func TestPool_UpdateTakeProfitLiabilitiesInvalid(t *testing.T) {
	// ctx := sdk.Context{} // mock or setup a context

	// Define the perpetual pool with assets
	ammPool := ammtypes.Pool{
		PoolId: 1,
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: sdk.NewInt(100),
				},
			},
		},
	}
	pool := types.NewPool(ammPool)
	pool.PoolAssetsLong = []types.PoolAsset{
		{
			Liabilities:           sdk.NewInt(0),
			TakeProfitLiabilities: sdk.NewInt(0),
			Custody:               sdk.NewInt(0),
			AssetDenom:            "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 50.
	denom := "testAsset2"
	err := pool.UpdateTakeProfitLiabilities(denom, sdk.NewInt(100), true, types.Position_LONG)
	// Expect that there is invalid asset denom error.
	assert.True(t, errors.Is(err, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")))

	// Expect that there is still 0 liabilities
	assert.Equal(t, pool.PoolAssetsLong[0].TakeProfitLiabilities, sdk.NewInt(0))
}

func TestPool_UpdateTakeProfitCustodyValid(t *testing.T) {

	// Define the perpetual pool with assets
	ammPool := ammtypes.Pool{
		PoolId: 1,
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: sdk.NewInt(100),
				},
			},
		},
	}
	pool := types.NewPool(ammPool)
	pool.PoolAssetsLong = []types.PoolAsset{
		{
			TakeProfitCustody: sdk.NewInt(0),
			Custody:           sdk.NewInt(0),
			AssetDenom:        "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 150.
	denom := "testAsset"
	err := pool.UpdateTakeProfitCustody(denom, sdk.NewInt(100), true, types.Position_LONG)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is 100 liabilities
	assert.Equal(t, pool.PoolAssetsLong[0].TakeProfitCustody, sdk.NewInt(100))
	err = pool.UpdateTakeProfitCustody(denom, sdk.NewInt(150), false, types.Position_LONG)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is -50 liabilities
	assert.Equal(t, pool.PoolAssetsLong[0].TakeProfitCustody, sdk.NewInt(-50))
}

func TestPool_UpdateTakeProfitCustodyInvalid(t *testing.T) {
	// ctx := sdk.Context{} // mock or setup a context

	// Define the perpetual pool with assets
	ammPool := ammtypes.Pool{
		PoolId: 1,
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: sdk.NewInt(100),
				},
			},
		},
	}
	pool := types.NewPool(ammPool)
	pool.PoolAssetsLong = []types.PoolAsset{
		{
			TakeProfitCustody: sdk.NewInt(0),
			Custody:           sdk.NewInt(0),
			AssetDenom:        "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 50.
	denom := "testAsset2"
	err := pool.UpdateTakeProfitCustody(denom, sdk.NewInt(100), true, types.Position_LONG)
	// Expect that there is invalid asset denom error.
	assert.True(t, errors.Is(err, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")))

	// Expect that there is still 0 liabilities
	assert.Equal(t, pool.PoolAssetsLong[0].TakeProfitCustody, sdk.NewInt(0))
}

func TestPool_UpdateCustodyValid(t *testing.T) {
	// ctx := sdk.Context{} // mock or setup a context

	// Define the perpetual pool with assets
	ammPool := ammtypes.Pool{
		PoolId: 1,
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: sdk.NewInt(100),
				},
			},
		},
	}
	pool := types.NewPool(ammPool)
	pool.PoolAssetsLong = []types.PoolAsset{
		{
			Liabilities: sdk.NewInt(0),
			Custody:     sdk.NewInt(0),
			AssetDenom:  "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 150.
	denom := "testAsset"
	err := pool.UpdateCustody(denom, sdk.NewInt(100), true, types.Position_LONG)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is 100 custody
	assert.Equal(t, pool.PoolAssetsLong[0].Custody, sdk.NewInt(100))
	err = pool.UpdateCustody(denom, sdk.NewInt(150), false, types.Position_LONG)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is -50 custody
	assert.Equal(t, pool.PoolAssetsLong[0].Custody, sdk.NewInt(-50))
}

func TestPool_UpdateCustodyInvalid(t *testing.T) {
	// ctx := sdk.Context{} // mock or setup a context

	// Define the perpetual pool with assets
	ammPool := ammtypes.Pool{
		PoolId: 1,
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: sdk.NewInt(100),
				},
			},
		},
	}
	pool := types.NewPool(ammPool)
	pool.PoolAssetsLong = []types.PoolAsset{
		{
			Liabilities: sdk.NewInt(0),
			Custody:     sdk.NewInt(0),
			AssetDenom:  "testAsset",
		},
	}

	// Test scenario, increase 100
	denom := "testAsset2"
	err := pool.UpdateCustody(denom, sdk.NewInt(100), true, types.Position_LONG)
	// Expect that there is invalid asset denom error.
	assert.True(t, errors.Is(err, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")))

	// Expect that there is still 0 custody
	assert.Equal(t, pool.PoolAssetsLong[0].Custody, sdk.NewInt(0))
}

func TestPool_InitiatePoolValid(t *testing.T) {
	// ctx := sdk.Context{} // mock or setup a context

	// Define the ammPool with assets
	ammPool := ammtypes.Pool{
		PoolId: 1,
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: sdk.NewInt(100),
				},
			},
		},
	}
	// Define the perpetual pool with assets
	pool := types.NewPool(ammPool)

	denom := "testAsset"
	assert.Equal(t, pool.AmmPoolId, (uint64)(1))
	assert.Equal(t, len(pool.PoolAssetsLong), 1)
	assert.Equal(t, len(pool.PoolAssetsShort), 1)
	assert.Equal(t, pool.PoolAssetsLong[0].AssetDenom, denom)
}
