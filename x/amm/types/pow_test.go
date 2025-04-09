package types_test

import (
	"testing"

	"github.com/elys-network/elys/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/stretchr/testify/require"
)

func TestPow(t *testing.T) {
	pow := types.Pow(osmomath.NewBigDec(2), osmomath.NewBigDecWithPrec(25, 1)) // 2^2.5
	require.Equal(t, pow.String(), "5.656854249492380195206754896838792316")
	pow = types.Pow(osmomath.NewBigDec(10), osmomath.NewBigDecWithPrec(25, 1)) // 10^2.5
	require.Equal(t, pow.String(), "316.227766016837933199889354443271853400")
}
