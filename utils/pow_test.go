package utils_test

import (
	"testing"

	"github.com/elys-network/elys/v4/utils"

	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/stretchr/testify/require"
)

func TestPow(t *testing.T) {
	pow := utils.Pow(osmomath.NewBigDec(2), osmomath.NewBigDecWithPrec(25, 1)) // 2^2.5
	require.Equal(t, pow.String(), "5.656854249492380196000000000000000000")
	pow = utils.Pow(osmomath.NewBigDec(10), osmomath.NewBigDecWithPrec(25, 1)) // 10^2.5
	require.Equal(t, pow.String(), "316.227766016837933200000000000000000000")
}
