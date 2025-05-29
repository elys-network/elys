package keeper

import (
	"github.com/elys-network/elys/v6/x/oracle/types"
)

var _ types.QueryServer = Keeper{}
