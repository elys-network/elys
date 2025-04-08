package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateEnableVestNow{}

func NewMsgUpdateAirdropParams(govAddress string, enableClaim bool, startAirdropHeight uint64, endAirdropHeight uint64, startKolHeight uint64, endKolHeight uint64) MsgUpdateAirdropParams {
	return MsgUpdateAirdropParams{
		Authority:               govAddress,
		EnableClaim:             enableClaim,
		StartAirdropClaimHeight: startAirdropHeight,
		EndAirdropClaimHeight:   endAirdropHeight,
		StartKolClaimHeight:     startKolHeight,
		EndKolClaimHeight:       endKolHeight,
	}

}

func (msg *MsgUpdateAirdropParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
