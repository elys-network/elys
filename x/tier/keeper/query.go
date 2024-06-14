package keeper

import (
	"github.com/elys-network/elys/x/tier/types"
)

var _ types.QueryServer = Keeper{}
