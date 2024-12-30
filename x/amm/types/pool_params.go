package types

import (
	sdkmath "cosmossdk.io/math"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	MaxSwapFee = sdkmath.LegacyMustNewDecFromStr("0.02") // 2%
)

func (params PoolParams) Validate() error {
	if params.SwapFee.IsNil() {
		return errors.New("swap_fee is nil")
	}

	if params.SwapFee.IsNegative() {
		return ErrNegativeSwapFee
	}

	if params.SwapFee.GT(MaxSwapFee) {
		return ErrTooMuchSwapFee
	}

	if err := sdk.ValidateDenom(params.FeeDenom); err != nil {
		return err
	}

	return nil
}
