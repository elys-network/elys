package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func NewPool(poolId uint64) Pool {
	return Pool{
		AmmPoolId:    poolId,
		Health:       sdk.NewDec(100),
		Enabled:      true,
		Closed:       false,
		InterestRate: sdk.NewDecFromIntWithPrec(sdk.NewInt(1), 1),
	}
}

// Update the asset liabilities
func (p *Pool) UpdateLiabilities(ctx sdk.Context, assetDenom string, amount sdk.Int, isIncrease bool) error {
	for i, asset := range p.PoolAssets {
		if asset.AssetDenom == assetDenom {
			if isIncrease {
				p.PoolAssets[i].Liabilities = asset.Liabilities.Add(amount)
			} else {
				p.PoolAssets[i].Liabilities = asset.Liabilities.Sub(amount)
			}

			return nil
		}
	}

	return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
}

// Update the asset custody
func (p *Pool) UpdateCustody(ctx sdk.Context, assetDenom string, amount sdk.Int, isIncrease bool) error {
	for i, asset := range p.PoolAssets {
		if asset.AssetDenom == assetDenom {
			if isIncrease {
				p.PoolAssets[i].Custody = asset.Custody.Add(amount)
			} else {
				p.PoolAssets[i].Custody = asset.Custody.Sub(amount)
			}
			return nil
		}
	}

	return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
}

// Initialite pool asset according to its corresponding amm pool assets.
func (p *Pool) InitiatePool(ctx sdk.Context, ammPool *ammtypes.Pool) error {
	if ammPool == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidType, "invalid amm pool")
	}

	// Set pool Id
	p.AmmPoolId = ammPool.PoolId
	return nil
}
