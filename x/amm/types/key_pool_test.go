package types_test

import (
	"encoding/binary"
	"testing"

	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestPoolKey(t *testing.T) {
	poolID := uint64(1234567890)

	expectedKey := make([]byte, 8)
	binary.BigEndian.PutUint64(expectedKey, poolID)
	expectedKey = append(expectedKey, []byte("/")...)

	resultKey := types.PoolKey(poolID)

	require.Equal(t, expectedKey, resultKey)
}
