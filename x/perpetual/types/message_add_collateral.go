package types

import (
	errorsmod "cosmossdk.io/errors"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddCollateral{}

func (msg *MsgAddCollateral) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err = msg.AddCollateral.Validate(); err != nil {
		return err
	}

	if msg.AddCollateral.IsZero() {
		return errors.New("add collateral cannot be zero")
	}

	if msg.Id == 0 {
		return errors.New("id cannot be zero")
	}

	if msg.PoolId == 0 {
		return errors.New("invalid pool id")
	}
	return nil
}
