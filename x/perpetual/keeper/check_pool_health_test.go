package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/perpetual/types/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCheckPoolHealth_PoolNotFound(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.PoolChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		PoolChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context

	poolId := uint64(1)

	// Mock behavior
	mockChecker.On("GetPool", ctx, poolId).Return(types.Pool{}, false)

	err := k.CheckPoolHealth(ctx, poolId)

	// Expect an error about invalid collateral asset
	assert.True(t, errors.Is(err, types.ErrInvalidBorrowingAsset))
}

func TestCheckPoolHealth_PoolDisabledOrClosed(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.PoolChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		PoolChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context

	poolId := uint64(1)
	pool := types.Pool{} // some mocked pool

	// Mock behavior
	mockChecker.On("GetPool", ctx, poolId).Return(pool, true)
	mockChecker.On("IsPoolEnabled", ctx, poolId).Return(false)

	err := k.CheckPoolHealth(ctx, poolId)

	// Expect an error about the pool being disabled or closed
	assert.True(t, errors.Is(err, types.ErrMTPDisabled))
}

func TestCheckPoolHealth_PoolHealthTooLow(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.PoolChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		PoolChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context

	poolId := uint64(1)
	pool := types.Pool{
		Health: sdk.NewDec(5), // mock a low health
		// ... other pool attributes
	}

	// Mock behavior
	mockChecker.On("GetPool", ctx, poolId).Return(pool, true)
	mockChecker.On("IsPoolEnabled", ctx, poolId).Return(true)
	mockChecker.On("IsPoolClosed", ctx, poolId).Return(false)
	mockChecker.On("GetPoolOpenThreshold", ctx).Return(sdk.NewDec(10)) // threshold higher than health

	err := k.CheckPoolHealth(ctx, poolId)

	// Expect an error about pool health being too low
	assert.True(t, errors.Is(err, types.ErrInvalidPosition))
}

func TestCheckPoolHealth_PoolIsHealthy(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.PoolChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		PoolChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context

	poolId := uint64(1)
	pool := types.Pool{
		Health: sdk.NewDec(15), // mock a good health
		// ... other pool attributes
	}

	// Mock behavior
	mockChecker.On("GetPool", ctx, poolId).Return(pool, true)
	mockChecker.On("IsPoolEnabled", ctx, poolId).Return(true)
	mockChecker.On("IsPoolClosed", ctx, poolId).Return(false)
	mockChecker.On("GetPoolOpenThreshold", ctx).Return(sdk.NewDec(10))

	err := k.CheckPoolHealth(ctx, poolId)

	// Expect no errors
	assert.Nil(t, err)
}
