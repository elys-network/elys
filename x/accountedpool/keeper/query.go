package keeper

import (
	"github.com/elys-network/elys/v7/x/accountedpool/types"
)

var _ types.QueryServer = Keeper{}
