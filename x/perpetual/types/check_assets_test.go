package types_test

import (
	"errors"
	"testing"

	errorsmod "cosmossdk.io/errors"
	"github.com/elys-network/elys/v4/x/perpetual/types"
	"github.com/stretchr/testify/assert"

	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
)

func TestCheckLongAssets_InvalidAssets(t *testing.T) {
	err := types.CheckLongAssets(ptypes.BaseCurrency, ptypes.BaseCurrency, ptypes.BaseCurrency)
	assert.True(t, errors.Is(err, errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid borrowing asset")))

	err = types.CheckLongAssets(ptypes.ATOM, ptypes.BaseCurrency, ptypes.BaseCurrency)
	assert.True(t, errors.Is(err, errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid borrowing asset")))
}

func TestCheckLongAssets_ValidAssets(t *testing.T) {
	err := types.CheckLongAssets(ptypes.BaseCurrency, ptypes.ATOM, ptypes.BaseCurrency)
	assert.Nil(t, err)

	err = types.CheckLongAssets(ptypes.ATOM, ptypes.ATOM, ptypes.BaseCurrency)
	assert.Nil(t, err)
}

func TestCheckShortAssets_InvalidAssets(t *testing.T) {
	// Test invalid cases for short positions
	err := types.CheckShortAssets(ptypes.ATOM, ptypes.BaseCurrency, ptypes.BaseCurrency)
	assert.True(t, errors.Is(err, errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "cannot short the base currency")))

	err = types.CheckShortAssets(ptypes.ATOM, ptypes.ATOM, ptypes.BaseCurrency)
	assert.True(t, errors.Is(err, errorsmod.Wrap(types.ErrInvalidCollateralAsset, "collateral asset cannot be the same as the borrowed asset in a short position")))

	err = types.CheckShortAssets(ptypes.ATOM, "btc", ptypes.BaseCurrency)
	assert.True(t, errors.Is(err, errorsmod.Wrap(types.ErrInvalidCollateralAsset, "collateral asset for a short position must be the base currency")))
}

func TestCheckShortAssets_ValidAssets(t *testing.T) {
	// Test valid case for short position
	err := types.CheckShortAssets(ptypes.BaseCurrency, ptypes.ATOM, ptypes.BaseCurrency)
	assert.Nil(t, err)
}
