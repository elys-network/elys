package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgBrokerAddCollateral{}

func NewMsgBrokerAddCollateral(creator string, amount math.Int, id int32, owner string) *MsgBrokerAddCollateral {
	return &MsgBrokerAddCollateral{
		Creator: creator,
		Amount:  amount,
		Id:      id,
		Owner:   owner,
	}
}

func (msg *MsgBrokerAddCollateral) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
