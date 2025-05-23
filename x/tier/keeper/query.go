package keeper

import (
	"github.com/elys-network/elys/v5/x/tier/types"
)

var _ types.QueryServer = Keeper{}
