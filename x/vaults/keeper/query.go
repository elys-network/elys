package keeper

import (
	"github.com/elys-network/elys/x/vaults/types"
)

var _ types.QueryServer = Keeper{}
