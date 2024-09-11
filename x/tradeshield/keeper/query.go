package keeper

import (
	"github.com/elys-network/elys/x/tradeshield/types"
)

var _ types.QueryServer = Keeper{}
