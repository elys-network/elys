package keeper

import (
	"github.com/elys-network/elys/v5/x/burner/types"
)

var _ types.QueryServer = Keeper{}
