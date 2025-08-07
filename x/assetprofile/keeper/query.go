package keeper

import (
	"github.com/elys-network/elys/v7/x/assetprofile/types"
)

var _ types.QueryServer = Keeper{}
