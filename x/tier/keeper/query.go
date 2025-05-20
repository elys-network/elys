package keeper

import (
	"github.com/elys-network/elys/v4/x/tier/types"
)

var _ types.QueryServer = Keeper{}
