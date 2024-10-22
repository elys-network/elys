package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func NewPool(poolId uint64, maxLeverage math.LegacyDec) Pool {
	return Pool{
		AmmPoolId:         poolId,
		Health:            math.LegacyOneDec(),
		LeveragedLpAmount: math.ZeroInt(),
		LeverageMax:       maxLeverage,
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
