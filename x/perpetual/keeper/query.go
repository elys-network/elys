package keeper

import (
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

var _ types.QueryServer = Keeper{}
