package types_test

import (
	"testing"

	"github.com/elys-network/elys/x/amm/types"
)

func TestGetPoolShareDenom(t *testing.T) {
	tests := []struct {
		poolId   uint64
		expected string
	}{
		{1, "amm/pool/1"},
		{42, "amm/pool/42"},
		{1000, "amm/pool/1000"},
	}

	for _, tt := range tests {
		result := types.GetPoolShareDenom(tt.poolId)
		if result != tt.expected {
			t.Errorf("GetPoolShareDenom(%d) = %s; want %s", tt.poolId, result, tt.expected)
		}
	}
}
