package keeper

import (
	"github.com/elys-network/elys/x/estaking/types"
)

var _ types.QueryServer = Keeper{}
