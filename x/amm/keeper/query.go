package keeper

import (
	"github.com/elys-network/elys/x/amm/types"
)

var _ types.QueryServer = Keeper{}
