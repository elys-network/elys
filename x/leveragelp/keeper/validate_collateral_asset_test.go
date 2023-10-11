package keeper_test

import (
	"errors"
	"testing"

	"github.com/elys-network/elys/x/leveragelp/keeper"
	"github.com/elys-network/elys/x/leveragelp/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateCollateralAsset_ValidCollateralAsset(t *testing.T) {
	k := keeper.Keeper{}

	collateralAsset := paramtypes.BaseCurrency // Correct asset

	err := k.ValidateCollateralAsset(collateralAsset)

	// Expect no error
	assert.Nil(t, err)
}

func TestValidateCollateralAsset_InvalidCollateralAsset(t *testing.T) {
	k := keeper.Keeper{}

	collateralAsset := "INVALID_ASSET" // Incorrect asset

	err := k.ValidateCollateralAsset(collateralAsset)

	// Expect an error about invalid collateral asset
	assert.True(t, errors.Is(err, types.ErrInvalidCollateralAsset))
}
