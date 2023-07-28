package keeper

import (
	"github.com/elys-network/elys/x/bridgeamm/types"
)

var _ types.QueryServer = Keeper{}
