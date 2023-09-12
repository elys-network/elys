package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCheckSameAssets_OpenPositionsBelowMax(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.PositionChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		PositionChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context

	// Mock behavior
	mockChecker.On("GetOpenMTPCount", ctx).Return(uint64(5))
	mockChecker.On("GetMaxOpenPositions", ctx).Return(uint64(10))

	err := k.CheckMaxOpenPositions(ctx)

	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}
