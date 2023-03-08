package keeper

import (
	"github.com/elys-network/elys/x/tokenomics/types"
)

var _ types.QueryServer = Keeper{}
