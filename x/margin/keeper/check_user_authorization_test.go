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

func TestCheckUserAuthorization_WhitelistingEnabledUserWhitelisted(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.AuthorizationChecker)

	// Create an instance of Keeper with the mock checker
	keeper := keeper.Keeper{
		AuthorizationChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context
	msg := &types.MsgOpen{Creator: "whitelistedUser"}

	// Mock behavior
	mockChecker.On("IsWhitelistingEnabled", ctx).Return(true)
	mockChecker.On("CheckIfWhitelisted", ctx, "whitelistedUser").Return(true)

	err := keeper.CheckUserAuthorization(ctx, msg)

	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}

func TestCheckUserAuthorization_WhitelistingEnabledUserNotWhitelisted(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.AuthorizationChecker)

	// Create an instance of Keeper with the mock checker
	keeper := keeper.Keeper{
		AuthorizationChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context
	msg := &types.MsgOpen{Creator: "nonWhitelistedUser"}

	// Mock behavior
	mockChecker.On("IsWhitelistingEnabled", ctx).Return(true)
	mockChecker.On("CheckIfWhitelisted", ctx, "nonWhitelistedUser").Return(false)

	err := keeper.CheckUserAuthorization(ctx, msg)

	// Expect an unauthorized error
	assert.True(t, errors.Is(err, types.ErrUnauthorised))
	mockChecker.AssertExpectations(t)
}

func TestCheckUserAuthorization_WhitelistingDisabled(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.AuthorizationChecker)

	// Create an instance of Keeper with the mock checker
	keeper := keeper.Keeper{
		AuthorizationChecker: mockChecker,
	}

	ctx := sdk.Context{}                      // mock or setup a context
	msg := &types.MsgOpen{Creator: "anyUser"} // Because whitelisting is off, user status doesn't matter.

	// Mock behavior
	mockChecker.On("IsWhitelistingEnabled", ctx).Return(false)

	err := keeper.CheckUserAuthorization(ctx, msg)

	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}
