package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/assert"
)

func TestCalcTakeAmount(t *testing.T) {
	custodyAmount := math.NewInt(int64(3000))
	fundingRate := math.LegacyMustNewDecFromStr("0.498826492562388470")
	got := CalcTakeAmount(custodyAmount, fundingRate)
	want := math.NewInt(int64(1496))
	assert.Equal(t, got, want)
}
