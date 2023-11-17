package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/elys-network/elys/x/margin/types/mocks"
	"github.com/stretchr/testify/assert"

	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func TestCheckShortAssets_InvalidAssets(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenShortChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context

	// Test invalid cases for short positions
	err := k.CheckShortAssets(ctx, ptypes.ATOM, ptypes.BaseCurrency, ptypes.BaseCurrency)
	assert.True(t, errors.Is(err, sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "cannot short the base currency")))

	err = k.CheckShortAssets(ctx, ptypes.ATOM, ptypes.ATOM, ptypes.BaseCurrency)
	assert.True(t, errors.Is(err, sdkerrors.Wrap(types.ErrInvalidCollateralAsset, "collateral asset cannot be the same as the borrowed asset in a short position")))

	err = k.CheckShortAssets(ctx, ptypes.ATOM, "btc", ptypes.BaseCurrency)
	assert.True(t, errors.Is(err, sdkerrors.Wrap(types.ErrInvalidCollateralAsset, "collateral asset for a short position must be the base currency")))

	// Expect no error
	mockChecker.AssertExpectations(t)
}

func TestCheckShortAssets_ValidAssets(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenShortChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context

	// Test valid case for short position
	err := k.CheckShortAssets(ctx, ptypes.BaseCurrency, ptypes.ATOM, ptypes.BaseCurrency)
	assert.Nil(t, err)

	// Expect no error
	mockChecker.AssertExpectations(t)
}
