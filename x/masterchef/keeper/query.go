package keeper

import (
	"github.com/elys-network/elys/x/masterchef/types"
)

var _ types.QueryServer = Keeper{}
