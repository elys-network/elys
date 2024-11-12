package types_test

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestPow(t *testing.T) {
	pow := types.Pow(sdkmath.LegacyNewDec(2), sdkmath.LegacyNewDecWithPrec(25, 1)) // 2^2.5
	require.Equal(t, pow.String(), "5.656854249492380196")
	pow = types.Pow(sdkmath.LegacyNewDec(10), sdkmath.LegacyNewDecWithPrec(25, 1)) // 10^2.5
	require.Equal(t, pow.String(), "316.227766016837933200")
}
