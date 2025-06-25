package keeper

import (
	"github.com/elys-network/elys/v6/x/tradeshield/types"
)

var _ types.QueryServer = Keeper{}
