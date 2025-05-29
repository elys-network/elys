package keeper

import (
	"github.com/elys-network/elys/v5/x/perpetual/types"
)

var _ types.QueryServer = Keeper{}
