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

func (pool Pool) GetTotalLongOpenInterest() math.Int {
	totalLongOpenInterest := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsLong {
		// We subtract asset.Collateral from totalCustodyLong because for long with collateral same as trading asset and user will
		// be charged for that the collateral as well even though they have already given that amount to the pool.
		// For LONG, asset.Custody will be 0 only for base currency but asset.Collateral won't be zero for base currency and trading asset
		// We subtract asset.Collateral only when asset is trading asset and in that case asset.Custody won't be zero
		// For base currency, asset.Collateral might not be 0 but asset.Custody will be 0 in LONG
		// !asset.Custody.IsZero() ensures that asset is trading asset for LONG
		if !asset.Custody.IsZero() {
			totalLongOpenInterest = totalLongOpenInterest.Add(asset.Custody).Sub(asset.Collateral)
		}
	}

	return totalLongOpenInterest
}

func (pool Pool) GetTotalShortOpenInterest() math.Int {
	totalShortOpenInterest := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsShort {
		totalShortOpenInterest = totalShortOpenInterest.Add(asset.Liabilities)
	}
	return totalShortOpenInterest
}

// GetNetOpenInterest calculates the net open interest for a given pool.
// Note: Net open interest should always be in terms of trading asset
func (pool Pool) GetNetOpenInterest() math.Int {
	totalLongOpenInterest := pool.GetTotalLongOpenInterest()
	totalShortOpenInterest := pool.GetTotalShortOpenInterest()

	// Net Open Interest = Long custody - Short Liabilities
	netOpenInterest := totalLongOpenInterest.Sub(totalShortOpenInterest)

	return netOpenInterest
}
