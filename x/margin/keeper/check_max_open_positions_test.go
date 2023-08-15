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

func TestCheckMaxOpenPositions_OpenPositionsBelowMax(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.PositionChecker)

	// Create an instance of Keeper with the mock checker
	keeper := keeper.Keeper{
		PositionChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context

	// Mock behavior
	mockChecker.On("GetOpenMTPCount", ctx).Return(uint64(5))
	mockChecker.On("GetMaxOpenPositions", ctx).Return(10)

	err := keeper.CheckMaxOpenPositions(ctx)

	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}

func TestCheckMaxOpenPositions_OpenPositionsEqualToMax(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.PositionChecker)

	// Create an instance of Keeper with the mock checker
	keeper := keeper.Keeper{
		PositionChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context

	// Mock behavior
	mockChecker.On("GetOpenMTPCount", ctx).Return(uint64(10))
	mockChecker.On("GetMaxOpenPositions", ctx).Return(10)

	err := keeper.CheckMaxOpenPositions(ctx)

	// Expect an error about max open positions
	assert.True(t, errors.Is(err, types.ErrMaxOpenPositions))
	mockChecker.AssertExpectations(t)
}

func TestCheckMaxOpenPositions_OpenPositionsExceedMax(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.PositionChecker)

	// Create an instance of Keeper with the mock checker
	keeper := keeper.Keeper{
		PositionChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context

	// Mock behavior
	mockChecker.On("GetOpenMTPCount", ctx).Return(uint64(11))
	mockChecker.On("GetMaxOpenPositions", ctx).Return(10)

	err := keeper.CheckMaxOpenPositions(ctx)

	// Expect an error about max open positions
	assert.True(t, errors.Is(err, types.ErrMaxOpenPositions))
	mockChecker.AssertExpectations(t)
}
