package keeper_test

import (
	"errors"
	"testing"

	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateCollateralAsset_ValidCollateralAsset(t *testing.T) {
	k := keeper.Keeper{}

	collateralAsset := paramtypes.BaseCurrency // Correct asset

	err := k.ValidateCollateralAsset(collateralAsset, paramtypes.BaseCurrency)

	// Expect no error
	assert.Nil(t, err)
}

func TestValidateCollateralAsset_InvalidCollateralAsset(t *testing.T) {
	k := keeper.Keeper{}

	collateralAsset := "INVALID_ASSET" // Incorrect asset

	err := k.ValidateCollateralAsset(collateralAsset, paramtypes.BaseCurrency)

	// Expect an error about invalid collateral asset
	assert.True(t, errors.Is(err, types.ErrInvalidCollateralAsset))
}
