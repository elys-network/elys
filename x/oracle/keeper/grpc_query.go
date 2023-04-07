package keeper

import (
	"github.com/elys-network/elys/x/oracle/types"
)

var _ types.QueryServer = Keeper{}
