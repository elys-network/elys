package types_test

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDebt(t *testing.T) {
	debt := types.Debt{
		Address:               sample.AccAddress(),
		Borrowed:              math.NewInt(100),
		InterestPaid:          math.NewInt(20),
		InterestStacked:       math.NewInt(30),
		BorrowTime:            10,
		LastInterestCalcTime:  10,
		LastInterestCalcBlock: 10,
	}

	require.Equal(t, debt.Address, debt.GetOwnerAccount().String())
	require.Equal(t, math.NewInt(110), debt.GetTotalLiablities())
}
