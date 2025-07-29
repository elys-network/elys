package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func NewPool(poolId uint64, maxLeverage, maxLeveragelpRatio sdkmath.LegacyDec) Pool {
	return Pool{
		AmmPoolId:          poolId,
		Health:             sdkmath.LegacyOneDec(),
		LeveragedLpAmount:  sdkmath.ZeroInt(),
		LeverageMax:        maxLeverage,
		MaxLeveragelpRatio: maxLeveragelpRatio,
	}
}

// Initialite pool asset according to its corresponding amm pool assets.
func (p *Pool) InitiatePool(ctx sdk.Context, ammPool *ammtypes.Pool) error {
	if ammPool == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidType, "invalid amm pool")
	}

	// Set pool Id
	p.AmmPoolId = ammPool.PoolId
	return nil
}

func (p *Pool) UpdateAssetLeveragedAmount(ctx sdk.Context, denom string, amount sdkmath.Int, isIncrease bool) {
	newAssetLevAmounts := make([]*AssetLeverageAmount, 0)
	for _, asset := range p.AssetLeverageAmounts {
		if asset.Denom == denom {
			if isIncrease {
				asset.LeveragedAmount = asset.LeveragedAmount.Add(amount)
			} else {
				asset.LeveragedAmount = asset.LeveragedAmount.Sub(amount)
			}
		}
		newAssetLevAmounts = append(newAssetLevAmounts, asset)
	}
	p.AssetLeverageAmounts = newAssetLevAmounts
}

func (p Pool) GetBigDecLeveragedLpAmount() osmomath.BigDec {
	return osmomath.BigDecFromSDKInt(p.LeveragedLpAmount)
}

func (p Pool) GetBigDecMaxLeveragelpRatio() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.MaxLeveragelpRatio)
}
