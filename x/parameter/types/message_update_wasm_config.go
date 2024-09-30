package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateWasmConfig{}

func NewMsgUpdateWasmConfig(creator string, wasmMaxLabelSize string, wasmMaxSize string, wasmMaxProposalWasmSize string) *MsgUpdateWasmConfig {
	return &MsgUpdateWasmConfig{
		Creator:                 creator,
		WasmMaxLabelSize:        wasmMaxLabelSize,
		WasmMaxSize:             wasmMaxSize,
		WasmMaxProposalWasmSize: wasmMaxProposalWasmSize,
	}
}

func (msg *MsgUpdateWasmConfig) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	res, ok := math.NewIntFromString(msg.WasmMaxLabelSize)

	if !ok || !res.IsPositive() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "wasm max label size must be a positive integer")
	}

	res, ok = math.NewIntFromString(msg.WasmMaxProposalWasmSize)

	if !ok || !res.IsPositive() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "wasm max proposal wasm size must be a positive integer")
	}

	res, ok = math.NewIntFromString(msg.WasmMaxSize)

	if !ok || !res.IsPositive() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "wasm max size must be a positive integer")
	}
	return nil
}
