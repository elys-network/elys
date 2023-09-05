package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/elys-network/elys/x/margin/types/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetFirstValidPool_NoPoolID(t *testing.T) {
	mockAmm := new(mocks.AmmKeeper)
	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil)

	ctx := sdk.Context{} // mock or setup a context
	borrowAsset := "testAsset"

	// Mock behavior
	mockAmm.On("GetAllPoolIdsWithDenom", ctx, borrowAsset).Return([]uint64{})

	_, err := k.GetFirstValidPool(ctx, borrowAsset)

	// Expect an error about invalid borrowing asset
	assert.True(t, errors.Is(err, types.ErrInvalidBorrowingAsset))
	mockAmm.AssertExpectations(t)
}

func TestGetFirstValidPool_ValidPoolID(t *testing.T) {
	mockAmm := new(mocks.AmmKeeper)
	k := keeper.NewKeeper(nil, nil, nil, "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga", mockAmm, nil, nil)

	ctx := sdk.Context{} // mock or setup a context
	borrowAsset := "testAsset"

	// Mock behavior
	mockAmm.On("GetAllPoolIdsWithDenom", ctx, borrowAsset).Return([]uint64{42, 43, 44})

	poolID, err := k.GetFirstValidPool(ctx, borrowAsset)

	// Expect no error and the first pool ID to be returned
	assert.Nil(t, err)
	assert.Equal(t, uint64(42), poolID)
	mockAmm.AssertExpectations(t)
}
