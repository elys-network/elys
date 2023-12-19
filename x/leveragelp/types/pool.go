package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func NewPool(poolId uint64) Pool {
	return Pool{
		AmmPoolId:         poolId,
		Health:            sdk.NewDec(100),
		Enabled:           true,
		Closed:            false,
		LeveragedLpAmount: sdk.ZeroInt(),
		LeverageMax:       sdk.NewDec(10),
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
