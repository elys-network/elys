package keeper

import (
	"github.com/elys-network/elys/x/launchpad/types"
)

var _ types.QueryServer = Keeper{}
