package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateWasmConfig = "update_wasm_config"

var _ sdk.Msg = &MsgUpdateWasmConfig{}

func NewMsgUpdateWasmConfig(creator string, wasmMaxLabelSize string, wasmMaxSize string, wasmMaxProposalWasmSize string) *MsgUpdateWasmConfig {
	return &MsgUpdateWasmConfig{
		Creator:                 creator,
		WasmMaxLabelSize:        wasmMaxLabelSize,
		WasmMaxSize:             wasmMaxSize,
		WasmMaxProposalWasmSize: wasmMaxProposalWasmSize,
	}
}

func (msg *MsgUpdateWasmConfig) Route() string {
	return RouterKey
}

func (msg *MsgUpdateWasmConfig) Type() string {
	return TypeMsgUpdateWasmConfig
}

func (msg *MsgUpdateWasmConfig) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateWasmConfig) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateWasmConfig) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	res, ok := sdk.NewIntFromString(msg.WasmMaxLabelSize)

	if !ok || !res.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "wasm max label size must be a positive integer")
	}

	res, ok = sdk.NewIntFromString(msg.WasmMaxProposalWasmSize)

	if !ok || !res.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "wasm max proposal wasm size must be a positive integer")
	}

	res, ok = sdk.NewIntFromString(msg.WasmMaxSize)

	if !ok || !res.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "wasm max size must be a positive integer")
	}
	return nil
}
