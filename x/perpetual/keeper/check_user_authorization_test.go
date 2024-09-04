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

func TestCheckUserAuthorization_WhitelistingEnabledUserWhitelisted(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.AuthorizationChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		AuthorizationChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context
	msg := &types.MsgOpen{Creator: "cosmos1xv9tklw7d82sezh9haa573wufgy59vmwe6xxe5"}

	// Mock behavior
	mockChecker.On("IsWhitelistingEnabled", ctx).Return(true)
	mockChecker.On("CheckIfWhitelisted", ctx, "whitelistedUser").Return(true)

	err := k.CheckUserAuthorization(ctx, msg)

	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}

func TestCheckUserAuthorization_WhitelistingEnabledUserNotWhitelisted(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.AuthorizationChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		AuthorizationChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context
	msg := &types.MsgOpen{Creator: "nonWhitelistedUser"}

	// Mock behavior
	mockChecker.On("IsWhitelistingEnabled", ctx).Return(true)
	mockChecker.On("CheckIfWhitelisted", ctx, "nonWhitelistedUser").Return(false)

	err := k.CheckUserAuthorization(ctx, msg)

	// Expect an unauthorized error
	assert.True(t, errors.Is(err, types.ErrUnauthorised))
	mockChecker.AssertExpectations(t)
}

func TestCheckUserAuthorization_WhitelistingDisabled(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.AuthorizationChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		AuthorizationChecker: mockChecker,
	}

	ctx := sdk.Context{}                      // mock or setup a context
	msg := &types.MsgOpen{Creator: "anyUser"} // Because whitelisting is off, user status doesn't matter.

	// Mock behavior
	mockChecker.On("IsWhitelistingEnabled", ctx).Return(false)

	err := k.CheckUserAuthorization(ctx, msg)

	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}
