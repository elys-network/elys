package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/assert"
)

type calcMTPTakeProfitCustodyTest struct {
	Description string
	Mtp         MTP
	Want        math.Int
}

func TestCalcMTPTakeProfitCustody(t *testing.T) {

	zero := math.LegacyMustNewDecFromStr("0")

	tests := []calcMTPTakeProfitCustodyTest{
		{
			Description: "take profit price is zero",
			Mtp: MTP{
				TakeProfitPrice: zero,
			},
			Want: math.ZeroInt(),
		},
		{
			Description: "LONG",
			Mtp: MTP{
				TakeProfitPrice: math.LegacyMustNewDecFromStr("1.23"),
				Position:        Position_LONG,
				Liabilities:     math.NewInt(int64(3000)),
			},
			Want: math.NewInt(int64(2439)),
		},
		{
			Description: "SHORT",
			Mtp: MTP{
				TakeProfitPrice: math.LegacyMustNewDecFromStr("0.98"),
				Position:        Position_SHORT,
				Liabilities:     math.NewInt(int64(3000)),
			},
			Want: math.NewInt(int64(2940)),
		},
	}

	t.Parallel()

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			got := CalcMTPTakeProfitCustody(test.Mtp)
			assert.Equal(t, test.Want, got)
		})
	}
}
