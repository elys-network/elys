package keeper

import (
	"github.com/elys-network/elys/v4/x/perpetual/types"
)

var _ types.QueryServer = Keeper{}
