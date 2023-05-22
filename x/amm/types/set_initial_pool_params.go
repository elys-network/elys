package types

import (
	"time"
)

// setInitialPoolParams
func (p *Pool) setInitialPoolParams(params PoolParams, sortedAssets []*PoolAsset, curBlockTime time.Time) error {
	p.PoolParams = &params

	return nil
}
