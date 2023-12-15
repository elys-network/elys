package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestPow(t *testing.T) {
	pow := types.Pow(sdk.NewDec(2), sdk.NewDecWithPrec(25, 1)) // 2^2.5
	require.Equal(t, pow.String(), "5.656854249492380196")
	pow = types.Pow(sdk.NewDec(10), sdk.NewDecWithPrec(25, 1)) // 10^2.5
	require.Equal(t, pow.String(), "316.227766016837933200")
}
