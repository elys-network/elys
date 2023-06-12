package types

import "fmt"

func GetPoolShareDenom(poolId uint64) string {
	return fmt.Sprintf("amm/pool/%d", poolId)
}
