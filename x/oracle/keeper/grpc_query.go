package keeper

import (
	"github.com/elys-network/elys/v5/x/oracle/types"
)

var _ types.QueryServer = Keeper{}
