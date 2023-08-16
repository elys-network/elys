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
	keeper := keeper.Keeper{}

	collateralAsset := paramtypes.USDC // Correct asset

	err := keeper.ValidateCollateralAsset(collateralAsset)

	// Expect no error
	assert.Nil(t, err)
}

func TestValidateCollateralAsset_InvalidCollateralAsset(t *testing.T) {
	keeper := keeper.Keeper{}

	collateralAsset := "INVALID_ASSET" // Incorrect asset

	err := keeper.ValidateCollateralAsset(collateralAsset)

	// Expect an error about invalid collateral asset
	assert.True(t, errors.Is(err, types.ErrInvalidCollateralAsset))
}
