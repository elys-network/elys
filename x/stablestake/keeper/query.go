package keeper

import (
	"github.com/elys-network/elys/x/stablestake/types"
)

var _ types.QueryServer = Keeper{}
