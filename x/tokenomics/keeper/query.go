package keeper

import (
	"github.com/elys-network/elys/v4/x/tokenomics/types"
)

var _ types.QueryServer = Keeper{}
