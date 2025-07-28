package keeper

import (
	"github.com/elys-network/elys/v7/x/tokenomics/types"
)

var _ types.QueryServer = Keeper{}
