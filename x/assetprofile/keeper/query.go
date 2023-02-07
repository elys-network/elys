package keeper

import (
	"github.com/elys-network/elys/x/assetprofile/types"
)

var _ types.QueryServer = Keeper{}
