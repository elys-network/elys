package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateEnableVestNow{}

func NewMsgUpdateAirdropParams(govAddress string, enableClaim bool, startHeight uint64, endHeight uint64) MsgUpdateAirdropParams {
	return MsgUpdateAirdropParams{
		Authority:               govAddress,
		EnableClaim:             enableClaim,
		StartAirdropClaimHeight: startHeight,
		EndAirdropClaimHeight:   endHeight,
	}

}

func (msg *MsgUpdateAirdropParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
