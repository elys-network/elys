package keeper

import (
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

var _ types.QueryServer = Keeper{}
