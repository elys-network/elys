package keeper

import (
	"github.com/elys-network/elys/v4/x/accountedpool/types"
)

var _ types.QueryServer = Keeper{}
