package keeper

import (
	"github.com/elys-network/elys/x/accountedpool/types"
)

var _ types.QueryServer = Keeper{}
