package keeper

import (
	"github.com/elys-network/elys/v4/x/burner/types"
)

var _ types.QueryServer = Keeper{}
