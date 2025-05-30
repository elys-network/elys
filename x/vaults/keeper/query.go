package keeper

import (
	"github.com/elys-network/elys/v5/x/vaults/types"
)

var _ types.QueryServer = Keeper{}
