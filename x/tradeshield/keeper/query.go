package keeper

import (
	"github.com/elys-network/elys/v4/x/tradeshield/types"
)

var _ types.QueryServer = Keeper{}
