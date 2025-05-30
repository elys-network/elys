package keeper

import (
	"github.com/elys-network/elys/v6/x/burner/types"
)

var _ types.QueryServer = Keeper{}
