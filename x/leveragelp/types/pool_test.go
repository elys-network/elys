package types_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/leveragelp/types"
	"github.com/stretchr/testify/assert"
)

func TestPool_InitiatePoolInvalid(t *testing.T) {
	ctx := sdk.Context{} // mock or setup a context

	pool := types.NewPool(1)
	err := pool.InitiatePool(ctx, nil)
	assert.True(t, errors.Is(err, sdkerrors.Wrap(sdkerrors.ErrInvalidType, "invalid amm pool")))
}
