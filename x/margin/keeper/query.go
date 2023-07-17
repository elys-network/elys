package keeper

import (
	"github.com/elys-network/elys/x/margin/types"
)

var _ types.QueryServer = Keeper{}
