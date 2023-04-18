package keeper

import (
	"github.com/elys-network/elys/x/incentive/types"
)

var _ types.QueryServer = Keeper{}
