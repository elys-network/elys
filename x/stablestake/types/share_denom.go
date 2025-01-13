package types

import (
	"strconv"
)

func GetShareDenomForPool(poolId uint64) string {
	if poolId == PoolId {
		return "stablestake/share"
	}
	return "stablestake/share/pool/" + strconv.FormatUint(poolId, 10)
}
