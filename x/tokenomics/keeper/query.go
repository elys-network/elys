package keeper

import (
	"github.com/elys-network/elys/v5/x/tokenomics/types"
)

var _ types.QueryServer = Keeper{}
