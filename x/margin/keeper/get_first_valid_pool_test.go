package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/elys-network/elys/x/margin/types/mocks"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
)

func TestGetFirstValidPool_NoPoolID(t *testing.T) {
	mockAmm := new(mocks.AmmKeeper)
	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil, nil, nil)

	ctx := sdk.Context{} // mock or setup a context
	collateralAsset := ptypes.BaseCurrency
	borrowAsset := "testAsset"
	denoms := []string{collateralAsset, borrowAsset}

	// Mock behavior
	mockAmm.On("GetPoolIdWithAllDenoms", ctx, denoms).Return(uint64(0), false)

	_, err := k.GetFirstValidPool(ctx, collateralAsset, borrowAsset)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
	mockAmm.AssertExpectations(t)
}

func TestGetFirstValidPool_ValidPoolID(t *testing.T) {
	mockAmm := new(mocks.AmmKeeper)
	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil, nil, nil)

	ctx := sdk.Context{} // mock or setup a context
	collateralAsset := ptypes.BaseCurrency
	borrowAsset := "testAsset"
	denoms := []string{collateralAsset, borrowAsset}

	// Mock behavior
	mockAmm.On("GetPoolIdWithAllDenoms", ctx, denoms).Return(uint64(1), true)

	poolID, err := k.GetFirstValidPool(ctx, collateralAsset, borrowAsset)

	// Expect no error and the first pool ID to be returned
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), poolID)
	mockAmm.AssertExpectations(t)
}
