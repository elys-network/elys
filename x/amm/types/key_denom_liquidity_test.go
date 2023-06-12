package types_test

import (
	"testing"

	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestDenomLiquidityKey(t *testing.T) {
	denom := "liquidityToken"

	expectedKey := []byte("liquidityToken/")

	resultKey := types.DenomLiquidityKey(denom)

	require.Equal(t, expectedKey, resultKey)
}
