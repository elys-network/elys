package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdatePoolParams{}

func NewMsgUpdatePoolParams(authority string, poolId uint64, poolParams PoolParams) *MsgUpdatePoolParams {
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

	if msg.PoolId == 0 {
		return fmt.Errorf("invalid pool id: %d", msg.PoolId)
	}

	if err = msg.PoolParams.Validate(); err != nil {
		return err
	}

	return nil
}
