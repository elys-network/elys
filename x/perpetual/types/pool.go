package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func NewPool(poolId uint64) Pool {
	return Pool{
		AmmPoolId:                            poolId,
		Health:                               sdk.NewDec(100),
		Enabled:                              true,
		Closed:                               false,
		BorrowInterestRate:                   sdk.NewDecFromIntWithPrec(sdk.NewInt(1), 1),
		PoolAssetsLong:                       []PoolAsset{},
		PoolAssetsShort:                      []PoolAsset{},
		LastHeightBorrowInterestRateComputed: 0,
		FundingRate:                          sdk.ZeroDec(),
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

// Get relevant pool asset based on position direction and asset denom
func (p *Pool) GetPoolAsset(position Position, assetDenom string) *PoolAsset {
	poolAssets := p.GetPoolAssets(position)
	for i, asset := range *poolAssets {
		if asset.AssetDenom == assetDenom {
			return &(*poolAssets)[i]
		}
	}
	return nil
}

// Update the asset balance
func (p *Pool) UpdateBalance(ctx sdk.Context, assetDenom string, amount math.Int, isIncrease bool, position Position) error {
	poolAsset := p.GetPoolAsset(position, assetDenom)
	if poolAsset == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
	}

	if isIncrease {
		poolAsset.AssetBalance = poolAsset.AssetBalance.Add(amount)
	} else {
		poolAsset.AssetBalance = poolAsset.AssetBalance.Sub(amount)
	}

	return nil
}

// Update the asset liabilities
func (p *Pool) UpdateLiabilities(ctx sdk.Context, assetDenom string, amount math.Int, isIncrease bool, position Position) error {
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

// Update the asset take profit liabilities
func (p *Pool) UpdateTakeProfitLiabilities(ctx sdk.Context, assetDenom string, amount math.Int, isIncrease bool, position Position) error {
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
func (p *Pool) UpdateTakeProfitCustody(ctx sdk.Context, assetDenom string, amount math.Int, isIncrease bool, position Position) error {
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
func (p *Pool) UpdateCustody(ctx sdk.Context, assetDenom string, amount math.Int, isIncrease bool, position Position) error {
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

// Update the unsettled liabilities balance
func (p *Pool) UpdateBlockBorrowInterest(ctx sdk.Context, assetDenom string, amount math.Int, isIncrease bool, position Position) error {
	poolAsset := p.GetPoolAsset(position, assetDenom)
	if poolAsset == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "invalid asset denom")
	}

	if isIncrease {
		poolAsset.BlockBorrowInterest = poolAsset.BlockBorrowInterest.Add(amount)
	} else {
		poolAsset.BlockBorrowInterest = poolAsset.BlockBorrowInterest.Sub(amount)
	}

	return nil
}

// Initialite pool asset according to its corresponding amm pool assets.
func (p *Pool) InitiatePool(ctx sdk.Context, ammPool *ammtypes.Pool) error {
	if ammPool == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidType, "invalid amm pool")
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
