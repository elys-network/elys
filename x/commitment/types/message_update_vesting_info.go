package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateVestingInfo = "update_vesting_info"

var _ sdk.Msg = &MsgUpdateVestingInfo{}

func NewMsgUpdateVestingInfo(creator string, baseDenom string, vestingDenom string, epochIdentifier string, numEpochs int64, vestNowFactor int64, numMaxVestings int64) *MsgUpdateVestingInfo {
	return &MsgUpdateVestingInfo{
		Authority:       creator,
		BaseDenom:       baseDenom,
		VestingDenom:    vestingDenom,
		EpochIdentifier: epochIdentifier,
		NumEpochs:       numEpochs,
		VestNowFactor:   vestNowFactor,
		NumMaxVestings:  numMaxVestings,
	}
}

func (msg *MsgUpdateVestingInfo) Route() string {
	return RouterKey
}

func (msg *MsgUpdateVestingInfo) Type() string {
	return TypeMsgUpdateVestingInfo
}

func (msg *MsgUpdateVestingInfo) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateVestingInfo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateVestingInfo) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
