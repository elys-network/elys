package keeper

import (
	"github.com/elys-network/elys/v5/x/tradeshield/types"
)

var _ types.QueryServer = Keeper{}
