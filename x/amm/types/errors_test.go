package types_test

import (
	"testing"

	"github.com/elys-network/elys/v6/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	require.Equal(t, types.ErrLimitMaxAmount.Error(), "calculated amount is larger than max amount")
	require.Equal(t, types.ErrLimitMinAmount.Error(), "calculated amount is less than min amount")
	require.Equal(t, types.ErrInvalidMathApprox.Error(), "invalid calculated result")
	require.Equal(t, types.ErrInvalidPool.Error(), "attempting to create an invalid pool")
	require.Equal(t, types.ErrDenomNotFoundInPool.Error(), "denom does not exist in pool")
	require.Equal(t, types.ErrEmptyRoutes.Error(), "routes not defined")
	require.Equal(t, types.ErrNegativeSwapFee.Error(), "swap fee is negative")
	require.Equal(t, types.ErrNegativeExitFee.Error(), "exit fee is negative")
	require.Equal(t, types.ErrTooMuchSwapFee.Error(), "swap fee should be less than 0.020000000000000000 (2.000000000000000000 %)")
	require.Equal(t, types.ErrTooManyTokensOut.Error(), "tx is trying to get more tokens out of the pool than exist")
	require.Equal(t, types.ErrInvalidPoolId.Error(), "invalid pool id")
}

func TestConstants(t *testing.T) {
	require.Equal(t, types.InvalidInputDenomsErrFormat, "input denoms must already exist in the pool (%s)")
	require.Equal(t, types.FormatRepeatingPoolAssetsNotAllowedErrFormat, "repeating pool assets not allowed, found %s")
	require.Equal(t, types.FormatNoPoolAssetFoundErrFormat, "can't find the PoolAsset (%s)")
	require.Equal(t, types.ErrMsgFormatSharesLargerThanMax, "%s resulted shares is larger than the max amount of %s")
}
