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

func TestCheckLongAssets_InvalidAssets(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenLongChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context

	err := k.CheckLongAssets(ctx, ptypes.BaseCurrency, ptypes.BaseCurrency, ptypes.BaseCurrency)
	assert.True(t, errors.Is(err, sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid borrowing asset")))

	err = k.CheckLongAssets(ctx, ptypes.ATOM, ptypes.BaseCurrency, ptypes.BaseCurrency)
	assert.True(t, errors.Is(err, sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid borrowing asset")))

	// Expect no error
	mockChecker.AssertExpectations(t)
}

func TestCheckLongAssets_ValidAssets(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenLongChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenLongChecker: mockChecker,
	}

	ctx := sdk.Context{} // mock or setup a context

	err := k.CheckLongAssets(ctx, ptypes.BaseCurrency, ptypes.ATOM, ptypes.BaseCurrency)
	assert.Nil(t, err)

	err = k.CheckLongAssets(ctx, ptypes.ATOM, ptypes.ATOM, ptypes.BaseCurrency)
	assert.Nil(t, err)

	// Expect an error about max open positions
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}
