package keeper

import (
	"github.com/elys-network/elys/x/liquidityprovider/types"
)

var _ types.QueryServer = Keeper{}
