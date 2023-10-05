package types_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/stretchr/testify/assert"
)

func TestPool_UpdateBalanceValid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

	// Define the leveragelp pool with assets
	pool := types.NewPool(1)
	pool.PoolAssets = []types.PoolAsset{
		{
			Liabilities:          sdk.NewInt(0),
			Custody:              sdk.NewInt(0),
			AssetBalance:         sdk.NewInt(0),
			UnsettledLiabilities: sdk.NewInt(0),
			BlockInterest:        sdk.NewInt(0),
			AssetDenom:           "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 50.
	denom := "testAsset"
	err := pool.UpdateBalance(ctx, denom, sdk.NewInt(100), true)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is 100 balance
	assert.Equal(t, pool.PoolAssets[0].AssetBalance, sdk.NewInt(100))
	err = pool.UpdateBalance(ctx, denom, sdk.NewInt(50), false)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is 100 balance
	assert.Equal(t, pool.PoolAssets[0].AssetBalance, sdk.NewInt(50))
}

func TestPool_UpdateBalanceInvalid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

	// Define the leveragelp pool with assets
	pool := types.NewPool(1)
	pool.PoolAssets = []types.PoolAsset{
		{
			Liabilities:          sdk.NewInt(0),
			Custody:              sdk.NewInt(0),
			AssetBalance:         sdk.NewInt(0),
			UnsettledLiabilities: sdk.NewInt(0),
			BlockInterest:        sdk.NewInt(0),
			AssetDenom:           "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 50.
	denom := "testAsset2"
	err := pool.UpdateBalance(ctx, denom, sdk.NewInt(100), true)
	// Expect that there is invalid asset denom error.
	assert.True(t, errors.Is(err, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")))

	// Expect that there is still 0 balance
	assert.Equal(t, pool.PoolAssets[0].AssetBalance, sdk.NewInt(0))
}

func TestPool_UpdateLiabilitiesValid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

	// Define the leveragelp pool with assets
	pool := types.NewPool(1)
	pool.PoolAssets = []types.PoolAsset{
		{
			Liabilities:          sdk.NewInt(0),
			Custody:              sdk.NewInt(0),
			AssetBalance:         sdk.NewInt(0),
			UnsettledLiabilities: sdk.NewInt(0),
			BlockInterest:        sdk.NewInt(0),
			AssetDenom:           "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 150.
	denom := "testAsset"
	err := pool.UpdateLiabilities(ctx, denom, sdk.NewInt(100), true)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is 100 liabilities
	assert.Equal(t, pool.PoolAssets[0].Liabilities, sdk.NewInt(100))
	err = pool.UpdateLiabilities(ctx, denom, sdk.NewInt(150), false)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is -50 liabilities
	assert.Equal(t, pool.PoolAssets[0].Liabilities, sdk.NewInt(-50))
}

func TestPool_UpdateLiabilitiesInvalid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

	// Define the leveragelp pool with assets
	pool := types.NewPool(1)
	pool.PoolAssets = []types.PoolAsset{
		{
			Liabilities:          sdk.NewInt(0),
			Custody:              sdk.NewInt(0),
			AssetBalance:         sdk.NewInt(0),
			UnsettledLiabilities: sdk.NewInt(0),
			BlockInterest:        sdk.NewInt(0),
			AssetDenom:           "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 50.
	denom := "testAsset2"
	err := pool.UpdateLiabilities(ctx, denom, sdk.NewInt(100), true)
	// Expect that there is invalid asset denom error.
	assert.True(t, errors.Is(err, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")))

	// Expect that there is still 0 liabilities
	assert.Equal(t, pool.PoolAssets[0].Liabilities, sdk.NewInt(0))
}

func TestPool_UpdateCustodyValid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

	// Define the leveragelp pool with assets
	pool := types.NewPool(1)
	pool.PoolAssets = []types.PoolAsset{
		{
			Liabilities:          sdk.NewInt(0),
			Custody:              sdk.NewInt(0),
			AssetBalance:         sdk.NewInt(0),
			UnsettledLiabilities: sdk.NewInt(0),
			BlockInterest:        sdk.NewInt(0),
			AssetDenom:           "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 150.
	denom := "testAsset"
	err := pool.UpdateCustody(ctx, denom, sdk.NewInt(100), true)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is 100 custody
	assert.Equal(t, pool.PoolAssets[0].Custody, sdk.NewInt(100))
	err = pool.UpdateCustody(ctx, denom, sdk.NewInt(150), false)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is -50 custody
	assert.Equal(t, pool.PoolAssets[0].Custody, sdk.NewInt(-50))
}

func TestPool_UpdateCustodyInvalid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

	// Define the leveragelp pool with assets
	pool := types.NewPool(1)
	pool.PoolAssets = []types.PoolAsset{
		{
			Liabilities:          sdk.NewInt(0),
			Custody:              sdk.NewInt(0),
			AssetBalance:         sdk.NewInt(0),
			UnsettledLiabilities: sdk.NewInt(0),
			BlockInterest:        sdk.NewInt(0),
			AssetDenom:           "testAsset",
		},
	}

	// Test scenario, increase 100
	denom := "testAsset2"
	err := pool.UpdateCustody(ctx, denom, sdk.NewInt(100), true)
	// Expect that there is invalid asset denom error.
	assert.True(t, errors.Is(err, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")))

	// Expect that there is still 0 custody
	assert.Equal(t, pool.PoolAssets[0].Custody, sdk.NewInt(0))
}

func TestPool_UpdateUnsettledLiabilitiesValid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

	// Define the leveragelp pool with assets
	pool := types.NewPool(1)
	pool.PoolAssets = []types.PoolAsset{
		{
			Liabilities:          sdk.NewInt(0),
			Custody:              sdk.NewInt(0),
			AssetBalance:         sdk.NewInt(0),
			UnsettledLiabilities: sdk.NewInt(0),
			BlockInterest:        sdk.NewInt(0),
			AssetDenom:           "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 150.
	denom := "testAsset"
	err := pool.UpdateUnsettledLiabilities(ctx, denom, sdk.NewInt(100), true)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is 100 UnsettledLiabilities
	assert.Equal(t, pool.PoolAssets[0].UnsettledLiabilities, sdk.NewInt(100))
	err = pool.UpdateUnsettledLiabilities(ctx, denom, sdk.NewInt(150), false)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is -50 UnsettledLiabilities
	assert.Equal(t, pool.PoolAssets[0].UnsettledLiabilities, sdk.NewInt(-50))
}

func TestPool_UpdateUnsettledLiabilitiesInvalid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

	// Define the leveragelp pool with assets
	pool := types.NewPool(1)
	pool.PoolAssets = []types.PoolAsset{
		{
			Liabilities:          sdk.NewInt(0),
			Custody:              sdk.NewInt(0),
			AssetBalance:         sdk.NewInt(0),
			UnsettledLiabilities: sdk.NewInt(0),
			BlockInterest:        sdk.NewInt(0),
			AssetDenom:           "testAsset",
		},
	}

	// Test scenario, increase 100
	denom := "testAsset2"
	err := pool.UpdateUnsettledLiabilities(ctx, denom, sdk.NewInt(100), true)
	// Expect that there is invalid asset denom error.
	assert.True(t, errors.Is(err, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")))

	// Expect that there is still 0 UnsettledLiabilities
	assert.Equal(t, pool.PoolAssets[0].UnsettledLiabilities, sdk.NewInt(0))
}

func TestPool_UpdateBlockInterestValid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

	// Define the leveragelp pool with assets
	pool := types.NewPool(1)
	pool.PoolAssets = []types.PoolAsset{
		{
			Liabilities:          sdk.NewInt(0),
			Custody:              sdk.NewInt(0),
			AssetBalance:         sdk.NewInt(0),
			UnsettledLiabilities: sdk.NewInt(0),
			BlockInterest:        sdk.NewInt(0),
			AssetDenom:           "testAsset",
		},
	}

	// Test scenario, increase 100 and decrease 150.
	denom := "testAsset"
	err := pool.UpdateBlockInterest(ctx, denom, sdk.NewInt(100), true)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is 100 BlockInterest
	assert.Equal(t, pool.PoolAssets[0].BlockInterest, sdk.NewInt(100))
	err = pool.UpdateBlockInterest(ctx, denom, sdk.NewInt(150), false)
	// Expect that there is no error
	assert.Nil(t, err)
	// Expect that there is -50 BlockInterest
	assert.Equal(t, pool.PoolAssets[0].BlockInterest, sdk.NewInt(-50))
}

func TestPool_UpdateBlockInterestInvalid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

	// Define the leveragelp pool with assets
	pool := types.NewPool(1)
	pool.PoolAssets = []types.PoolAsset{
		{
			Liabilities:          sdk.NewInt(0),
			Custody:              sdk.NewInt(0),
			AssetBalance:         sdk.NewInt(0),
			UnsettledLiabilities: sdk.NewInt(0),
			BlockInterest:        sdk.NewInt(0),
			AssetDenom:           "testAsset",
		},
	}

	// Test scenario, increase 100
	denom := "testAsset2"
	err := pool.UpdateBlockInterest(ctx, denom, sdk.NewInt(100), true)
	// Expect that there is invalid asset denom error.
	assert.True(t, errors.Is(err, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")))

	// Expect that there is still 0 BlockInterest
	assert.Equal(t, pool.PoolAssets[0].BlockInterest, sdk.NewInt(0))
}

func TestPool_InitiatePoolValid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

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
	// Define the leveragelp pool with assets
	pool := types.NewPool(1)
	err := pool.InitiatePool(ctx, &ammPool)

	// Expect that there is no error
	assert.Nil(t, err)

	denom := "testAsset"
	assert.Equal(t, pool.AmmPoolId, (uint64)(1))
	assert.Equal(t, len(pool.PoolAssets), 1)
	assert.Equal(t, pool.PoolAssets[0].AssetDenom, denom)
}

func TestPool_InitiatePoolInvalid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

	pool := types.NewPool(1)
	err := pool.InitiatePool(ctx, nil)
	assert.True(t, errors.Is(err, sdkerrors.Wrap(sdkerrors.ErrInvalidType, "invalid amm pool")))
}
