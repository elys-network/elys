package types_test

import (
	"errors"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/stretchr/testify/assert"

	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func TestCheckShortAssets_InvalidAssets(t *testing.T) {
	// Test invalid cases for short positions
	err := types.CheckShortAssets(ptypes.ATOM, ptypes.BaseCurrency, ptypes.BaseCurrency)
	assert.True(t, errors.Is(err, sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "cannot short the base currency")))

	err = types.CheckShortAssets(ptypes.ATOM, ptypes.ATOM, ptypes.BaseCurrency)
	assert.True(t, errors.Is(err, sdkerrors.Wrap(types.ErrInvalidCollateralAsset, "collateral asset cannot be the same as the borrowed asset in a short position")))

	err = types.CheckShortAssets(ptypes.ATOM, "btc", ptypes.BaseCurrency)
	assert.True(t, errors.Is(err, sdkerrors.Wrap(types.ErrInvalidCollateralAsset, "collateral asset for a short position must be the base currency")))
}

func TestCheckShortAssets_ValidAssets(t *testing.T) {
	// Test valid case for short position
	err := types.CheckShortAssets(ptypes.BaseCurrency, ptypes.ATOM, ptypes.BaseCurrency)
	assert.Nil(t, err)
}
