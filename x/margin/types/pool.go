package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func NewPool(poolId uint64) Pool {
	return Pool{
		AmmPoolId:          poolId,
		Health:             sdk.NewDec(100),
		Enabled:            true,
		Closed:             false,
		BorrowInterestRate: sdk.NewDecFromIntWithPrec(sdk.NewInt(1), 1),
		PoolAssetsLong:     []PoolAsset{},
		PoolAssetsShort:    []PoolAsset{},
	}
}

// Get relevant pool asset array based on position direction
func (p *Pool) GetPoolAssets(position Position) *[]PoolAsset {
	if position == Position_LONG {
		return &p.PoolAssetsLong
	} else {
		return &p.PoolAssetsShort
	}
}

// Update the asset balance
func (p *Pool) UpdateBalance(ctx sdk.Context, assetDenom string, amount sdk.Int, isIncrease bool, position Position) error {
	poolAssets := p.GetPoolAssets(position)
	for i, asset := range *poolAssets {
		if asset.AssetDenom == assetDenom {
			if isIncrease {
				(*poolAssets)[i].AssetBalance = asset.AssetBalance.Add(amount)
			} else {
				(*poolAssets)[i].AssetBalance = asset.AssetBalance.Sub(amount)
			}

			return nil
		}
	}

	return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
}

// Update the asset liabilities
func (p *Pool) UpdateLiabilities(ctx sdk.Context, assetDenom string, amount sdk.Int, isIncrease bool, position Position) error {
	poolAssets := p.GetPoolAssets(position)
	for i, asset := range *poolAssets {
		if asset.AssetDenom == assetDenom {
			if isIncrease {
				(*poolAssets)[i].Liabilities = asset.Liabilities.Add(amount)
			} else {
				(*poolAssets)[i].Liabilities = asset.Liabilities.Sub(amount)
			}

			return nil
		}
	}

	return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
}

// Update the asset take profit liabilities
func (p *Pool) UpdateTakeProfitLiabilities(ctx sdk.Context, assetDenom string, amount sdk.Int, isIncrease bool, position Position) error {
	poolAssets := p.GetPoolAssets(position)
	for i, asset := range *poolAssets {
		if asset.AssetDenom == assetDenom {
			if isIncrease {
				(*poolAssets)[i].TakeProfitLiabilities = asset.TakeProfitLiabilities.Add(amount)
			} else {
				(*poolAssets)[i].TakeProfitLiabilities = asset.TakeProfitLiabilities.Sub(amount)
			}

			return nil
		}
	}

	return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
}

// Update the asset take profit custody
func (p *Pool) UpdateTakeProfitCustody(ctx sdk.Context, assetDenom string, amount sdk.Int, isIncrease bool, position Position) error {
	poolAssets := p.GetPoolAssets(position)
	for i, asset := range *poolAssets {
		if asset.AssetDenom == assetDenom {
			if isIncrease {
				(*poolAssets)[i].TakeProfitCustody = asset.TakeProfitCustody.Add(amount)
			} else {
				(*poolAssets)[i].TakeProfitCustody = asset.TakeProfitCustody.Sub(amount)
			}

			return nil
		}
	}

	return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
}

// Update the asset custody
func (p *Pool) UpdateCustody(ctx sdk.Context, assetDenom string, amount sdk.Int, isIncrease bool, position Position) error {
	poolAssets := p.GetPoolAssets(position)
	for i, asset := range *poolAssets {
		if asset.AssetDenom == assetDenom {
			if isIncrease {
				(*poolAssets)[i].Custody = asset.Custody.Add(amount)
			} else {
				(*poolAssets)[i].Custody = asset.Custody.Sub(amount)
			}
			return nil
		}
	}

	return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
}

// Update the unsettled liabilities balance
func (p *Pool) UpdateBlockBorrowInterest(ctx sdk.Context, assetDenom string, amount sdk.Int, isIncrease bool, position Position) error {
	poolAssets := p.GetPoolAssets(position)
	for i, asset := range *poolAssets {
		if asset.AssetDenom == assetDenom {
			if isIncrease {
				(*poolAssets)[i].BlockBorrowInterest = asset.BlockBorrowInterest.Add(amount)
			} else {
				(*poolAssets)[i].BlockBorrowInterest = asset.BlockBorrowInterest.Sub(amount)
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

	for _, asset := range ammPool.PoolAssets {
		poolAsset := PoolAsset{
			Liabilities:         sdk.ZeroInt(),
			Custody:             sdk.ZeroInt(),
			AssetBalance:        sdk.ZeroInt(),
			BlockBorrowInterest: sdk.ZeroInt(),
			AssetDenom:          asset.Token.Denom,
		}

		p.PoolAssetsLong = append(p.PoolAssetsLong, poolAsset)
		p.PoolAssetsShort = append(p.PoolAssetsShort, poolAsset)
	}

	return nil
}
