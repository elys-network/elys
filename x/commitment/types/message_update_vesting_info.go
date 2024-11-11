package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateVestingInfo{}

func NewMsgUpdateVestingInfo(creator string, baseDenom string, vestingDenom string, epochIdentifier string, numBlocks int64, vestNowFactor int64, numMaxVestings int64) *MsgUpdateVestingInfo {
	return &MsgUpdateVestingInfo{
		Authority:      creator,
		BaseDenom:      baseDenom,
		VestingDenom:   vestingDenom,
		NumBlocks:      numBlocks,
		VestNowFactor:  vestNowFactor,
		NumMaxVestings: numMaxVestings,
	}
}

func (msg *MsgUpdateVestingInfo) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
