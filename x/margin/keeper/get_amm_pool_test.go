package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/elys-network/elys/x/margin/types/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetAmmPool_PoolNotFound(t *testing.T) {
	mockAmm := new(mocks.AmmKeeper)
	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil, nil)

	ctx := sdk.Context{} // mock or setup a context
	borrowAsset := "testAsset"
	poolId := uint64(42)

	// Mock behavior
	mockAmm.On("GetPool", ctx, poolId).Return(ammtypes.Pool{}, false)

	_, err := k.GetAmmPool(ctx, poolId, borrowAsset)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
	mockAmm.AssertExpectations(t)
}

func TestGetAmmPool_PoolFound(t *testing.T) {
	mockAmm := new(mocks.AmmKeeper)
	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil, nil)

	ctx := sdk.Context{} // mock or setup a context
	borrowAsset := "testAsset"
	poolId := uint64(42)

	expectedPool := ammtypes.Pool{}

	// Mock behavior
	mockAmm.On("GetPool", ctx, poolId).Return(expectedPool, true)

	pool, err := k.GetAmmPool(ctx, poolId, borrowAsset)

	// Expect no error and the correct pool to be returned
	assert.Nil(t, err)
	assert.Equal(t, expectedPool, pool)
	mockAmm.AssertExpectations(t)
}
