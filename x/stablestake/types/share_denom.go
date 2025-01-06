package types

import "strconv"

func GetShareDenom() string {
	return "stablestake/share"
}

func GetShareDenomForPool(poolId uint64) string {
	return "stablestake/share/pool/" + strconv.FormatUint(poolId, 10)
}
