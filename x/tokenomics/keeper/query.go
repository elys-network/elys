package keeper

import (
	"github.com/elys-network/elys/v6/x/tokenomics/types"
)

var _ types.QueryServer = Keeper{}
