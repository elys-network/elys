package keeper

import (
	"github.com/elys-network/elys/v4/x/assetprofile/types"
)

var _ types.QueryServer = Keeper{}
