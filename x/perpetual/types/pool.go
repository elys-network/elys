package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func NewPool(ammPool ammtypes.Pool) Pool {
	p := Pool{
		AmmPoolId:                            ammPool.PoolId,
		Health:                               math.LegacyOneDec(),
		BorrowInterestRate:                   math.LegacyZeroDec(),
		PoolAssetsLong:                       []PoolAsset{},
		PoolAssetsShort:                      []PoolAsset{},
		LastHeightBorrowInterestRateComputed: 0,
		FundingRate:                          math.LegacyZeroDec(),
	}

	for _, asset := range ammPool.PoolAssets {
		poolAsset := PoolAsset{
			Liabilities: math.ZeroInt(),
			Custody:     math.ZeroInt(),
			AssetDenom:  asset.Token.Denom,
		}

		p.PoolAssetsLong = append(p.PoolAssetsLong, poolAsset)
		p.PoolAssetsShort = append(p.PoolAssetsShort, poolAsset)
	}

	return p
}

// Get relevant pool asset array based on position direction
func (p Pool) GetPoolAssets(position Position) *[]PoolAsset {
	if position == Position_LONG {
		return &p.PoolAssetsLong
	} else {
		return &p.PoolAssetsShort
	}
}

// Get relevant pool asset based on position direction and asset denom
func (p Pool) GetPoolAsset(position Position, assetDenom string) *PoolAsset {
	poolAssets := p.GetPoolAssets(position)
	for i, asset := range *poolAssets {
		if asset.AssetDenom == assetDenom {
			return &(*poolAssets)[i]
		}
	}
	return nil
}

// Update the asset liabilities
func (p *Pool) UpdateLiabilities(assetDenom string, amount math.Int, isIncrease bool, position Position) error {
	poolAsset := p.GetPoolAsset(position, assetDenom)
	if poolAsset == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
	}

	if isIncrease {
		poolAsset.Liabilities = poolAsset.Liabilities.Add(amount)
	} else {
		poolAsset.Liabilities = poolAsset.Liabilities.Sub(amount)
	}

	return nil
}

// Update the asset collateral
func (p *Pool) UpdateCollateral(assetDenom string, amount math.Int, isIncrease bool, position Position) error {
	poolAsset := p.GetPoolAsset(position, assetDenom)
	if poolAsset == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
	}

	if isIncrease {
		poolAsset.Collateral = poolAsset.Collateral.Add(amount)
	} else {
		poolAsset.Collateral = poolAsset.Collateral.Sub(amount)
	}

	return nil
}

// Update the asset take profit liabilities
func (p *Pool) UpdateTakeProfitLiabilities(assetDenom string, amount math.Int, isIncrease bool, position Position) error {
	poolAsset := p.GetPoolAsset(position, assetDenom)
	if poolAsset == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
	}

	if isIncrease {
		poolAsset.TakeProfitLiabilities = poolAsset.TakeProfitLiabilities.Add(amount)
	} else {
		poolAsset.TakeProfitLiabilities = poolAsset.TakeProfitLiabilities.Sub(amount)
	}

	return nil
}

// Update the asset take profit custody
func (p *Pool) UpdateTakeProfitCustody(assetDenom string, amount math.Int, isIncrease bool, position Position) error {
	poolAsset := p.GetPoolAsset(position, assetDenom)
	if poolAsset == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
	}

	if isIncrease {
		poolAsset.TakeProfitCustody = poolAsset.TakeProfitCustody.Add(amount)
	} else {
		poolAsset.TakeProfitCustody = poolAsset.TakeProfitCustody.Sub(amount)
	}

	return nil
}

// Update the asset custody
func (p *Pool) UpdateCustody(assetDenom string, amount math.Int, isIncrease bool, position Position) error {
	poolAsset := p.GetPoolAsset(position, assetDenom)
	if poolAsset == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
	}

	if isIncrease {
		poolAsset.Custody = poolAsset.Custody.Add(amount)
	} else {
		poolAsset.Custody = poolAsset.Custody.Sub(amount)
	}

	return nil
}

// Update the fees collected
func (p *Pool) UpdateFeesCollected(assetDenom string, amount math.Int, isIncrease bool) error {
	if isIncrease {
		for _, coin := range p.FeesCollected {
			if coin.Denom == assetDenom {
				coin.Amount = coin.Amount.Add(amount)
				return nil
			}
		}
	} else {
		for _, coin := range p.FeesCollected {
			if coin.Denom == assetDenom {
				coin.Amount = coin.Amount.Sub(amount)
				return nil
			}
		}
	}
	p.FeesCollected = append(p.FeesCollected, sdk.NewCoin(assetDenom, amount))

	return nil
}
