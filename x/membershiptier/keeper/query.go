package keeper

import (
	"github.com/elys-network/elys/x/membershiptier/types"
)

var _ types.QueryServer = Keeper{}
