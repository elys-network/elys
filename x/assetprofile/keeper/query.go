package keeper

import (
	"github.com/elys-network/elys/v5/x/assetprofile/types"
)

var _ types.QueryServer = Keeper{}
