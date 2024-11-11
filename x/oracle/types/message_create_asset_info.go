package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateAssetInfo{}

func NewMsgCreateAssetInfo(creator string, denom string, display string, bandTicker string, elysTicker string, decimal uint64) *MsgCreateAssetInfo {
	return &MsgCreateAssetInfo{
		Creator:    creator,
		Denom:      denom,
		Display:    display,
		BandTicker: bandTicker,
		ElysTicker: elysTicker,
		Decimal:    decimal,
	}
}

func (msg *MsgCreateAssetInfo) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
