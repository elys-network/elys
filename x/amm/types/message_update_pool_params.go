package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdatePoolParams{}

func NewMsgUpdatePoolParams(authority string, poolId uint64, poolParams *PoolParams) *MsgUpdatePoolParams {
	return &MsgUpdatePoolParams{
		Authority:  authority,
		PoolId:     poolId,
		PoolParams: poolParams,
	}
}

func (msg *MsgUpdatePoolParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.PoolParams == nil {
		return ErrPoolParamsShouldNotBeNil
	}

	if msg.PoolParams.SwapFee.IsNegative() {
		return ErrFeeShouldNotBeNegative
	}

	if msg.PoolParams.SwapFee.GT(sdkmath.LegacyNewDecWithPrec(2, 2)) { // >2%
		return ErrSwapFeeShouldNotExceedTwoPercent
	}

	if msg.PoolParams.ExitFee.IsNegative() {
		return ErrFeeShouldNotBeNegative
	}

	if msg.PoolParams.ExitFee.GT(sdkmath.LegacyNewDecWithPrec(2, 2)) { // >2%
		return ErrExitFeeShouldNotExceedTwoPercent
	}

	return nil
}
