package keeper

import (
	"github.com/elys-network/elys/v7/x/oracle/types"
)

var _ types.QueryServer = Keeper{}
