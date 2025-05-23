package keeper

import (
	"github.com/elys-network/elys/v5/x/accountedpool/types"
)

var _ types.QueryServer = Keeper{}
