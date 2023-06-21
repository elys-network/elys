package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetInitialPoolAssets sets the PoolAssets in the pool. It is only designed to
// be called at the pool's creation. If the same denom's PoolAsset exists, will
// return error.
//
// The list of PoolAssets must be sorted. This is done to enable fast searching
// for a PoolAsset by denomination.
// TODO: Unify story for validation of []PoolAsset, some is here, some is in
// CreatePool.ValidateBasic()
func (p *Pool) SetInitialPoolAssets(PoolAssets []PoolAsset) error {
	exists := make(map[string]bool)
	for _, asset := range p.PoolAssets {
		exists[asset.Token.Denom] = true
	}

	newTotalWeight := p.TotalWeight
	scaledPoolAssets := make([]PoolAsset, 0, len(PoolAssets))

	// TODO: Refactor this into PoolAsset.validate()
	for _, asset := range PoolAssets {
		if asset.Token.Amount.LTE(sdk.ZeroInt()) {
			return fmt.Errorf("can't add the zero or negative balance of token")
		}

		err := asset.validateWeight()
		if err != nil {
			return err
		}

		if exists[asset.Token.Denom] {
			return fmt.Errorf("same PoolAsset already exists")
		}
		exists[asset.Token.Denom] = true

		// Scale weight from the user provided weight to the correct internal weight
		asset.Weight = asset.Weight.MulRaw(GuaranteedWeightPrecision)
		scaledPoolAssets = append(scaledPoolAssets, asset)
		newTotalWeight = newTotalWeight.Add(asset.Weight)
	}

	// TODO: Change this to a more efficient sorted insert algorithm.
	// Furthermore, consider changing the underlying data type to allow in-place modification if the
	// number of PoolAssets is expected to be large.
	p.PoolAssets = append(p.PoolAssets, scaledPoolAssets...)

	sortPoolAssetsByDenom(p.PoolAssets)

	p.TotalWeight = newTotalWeight

	return nil
}
