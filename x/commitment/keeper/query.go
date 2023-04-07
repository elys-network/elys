package keeper

import (
	"github.com/elys-network/elys/x/commitment/types"
)

var _ types.QueryServer = Keeper{}
