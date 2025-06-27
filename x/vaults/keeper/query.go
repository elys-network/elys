package keeper

import (
	"github.com/elys-network/elys/v6/x/vaults/types"
)

var _ types.QueryServer = Keeper{}
