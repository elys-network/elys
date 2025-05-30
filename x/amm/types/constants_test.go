package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/elys-network/elys/v6/x/amm/types"
)

func TestOneShare(t *testing.T) {
	expected := math.NewIntWithDecimal(1, types.OneShareExponent)
	if !types.OneShare.Equal(expected) {
		t.Errorf("OneShare = %s; want %s", types.OneShare.String(), expected.String())
	}
}

func TestInitPoolSharesSupply(t *testing.T) {
	expected := types.OneShare.MulRaw(100)
	if !types.InitPoolSharesSupply.Equal(expected) {
		t.Errorf("InitPoolSharesSupply = %s; want %s", types.InitPoolSharesSupply.String(), expected.String())
	}
}

func TestGuaranteedWeightPrecision(t *testing.T) {
	expected := int64(1 << 30)
	if types.GuaranteedWeightPrecision != expected {
		t.Errorf("GuaranteedWeightPrecision = %d; want %d", types.GuaranteedWeightPrecision, expected)
	}
}
