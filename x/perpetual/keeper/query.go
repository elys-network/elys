package keeper

import (
	"github.com/elys-network/elys/x/perpetual/types"
)

var _ types.QueryServer = Keeper{}
