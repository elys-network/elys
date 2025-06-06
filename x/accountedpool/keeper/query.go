package keeper

import (
	"github.com/elys-network/elys/v6/x/accountedpool/types"
)

var _ types.QueryServer = Keeper{}
