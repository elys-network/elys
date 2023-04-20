package keeper

import (
	"github.com/elys-network/elys/x/burner/types"
)

var _ types.QueryServer = Keeper{}
