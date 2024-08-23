package types_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/stretchr/testify/assert"
)

func TestPool_InitiatePoolInvalid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

	pool := types.NewPool(1)
	err := pool.InitiatePool(ctx, nil)
	assert.True(t, errors.Is(err, errorsmod.Wrap(sdkerrors.ErrInvalidType, "invalid amm pool")))

	err = pool.InitiatePool(ctx, &ammtypes.Pool{PoolId: 1})
	require.NoError(t, err)
}
