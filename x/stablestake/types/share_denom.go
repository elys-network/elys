package types

import (
	fmt "fmt"
	"strconv"
	"strings"
)

func GetShareDenomForPool(poolId uint64) string {
	if poolId == UsdcPoolId {
		return "stablestake/share"
	}
	return "stablestake/share/pool/" + strconv.FormatUint(poolId, 10)
}

// GetPoolIDFromPath retrieves the poolid from the given path in the format "stablestake/share/pool/poolid".
func GetPoolIDFromPath(path string) (uint64, error) {
	if path == "stablestake/share" {
		return UsdcPoolId, nil
	}
	parts := strings.Split(path, "/")
	if len(parts) != 4 || parts[0] != "stablestake" || parts[1] != "share" || parts[2] != "pool" {
		return 0, fmt.Errorf("invalid path format")
	}
	poolID, err := strconv.ParseUint(parts[3], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid pool ID: %v", err)
	}
	return poolID, nil
}
