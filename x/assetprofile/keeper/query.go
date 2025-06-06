package keeper

import (
	"github.com/elys-network/elys/v6/x/assetprofile/types"
)

var _ types.QueryServer = Keeper{}
