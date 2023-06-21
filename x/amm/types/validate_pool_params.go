package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (params PoolParams) Validate(poolWeights []PoolAsset) error {
	if params.ExitFee.IsNegative() {
		return ErrNegativeExitFee
	}

	if params.ExitFee.GTE(sdk.OneDec()) {
		return ErrTooMuchExitFee
	}

	if params.SwapFee.IsNegative() {
		return ErrNegativeSwapFee
	}

	if params.SwapFee.GTE(sdk.OneDec()) {
		return ErrTooMuchSwapFee
	}

	return nil
}
