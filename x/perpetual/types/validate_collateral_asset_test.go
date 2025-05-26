package types_test

import (
	"errors"
	"testing"

	paramtypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/elys-network/elys/v5/x/perpetual/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateCollateralAsset_ValidCollateralAsset(t *testing.T) {
	collateralAsset := paramtypes.BaseCurrency // Correct asset

	err := types.ValidateCollateralAsset(collateralAsset, paramtypes.BaseCurrency)

	// Expect no error
	assert.Nil(t, err)
}

func TestValidateCollateralAsset_InvalidCollateralAsset(t *testing.T) {
	collateralAsset := "INVALID_ASSET" // Incorrect asset

	err := types.ValidateCollateralAsset(collateralAsset, paramtypes.BaseCurrency)

	// Expect an error about invalid collateral asset
	assert.True(t, errors.Is(err, types.ErrInvalidCollateralAsset))
}
