package keeper

import (
	"github.com/elys-network/elys/v4/x/oracle/types"
)

var _ types.QueryServer = Keeper{}
