package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateEntry{}

func NewMsgUpdateEntry(
	authority string,
	baseDenom string,
	decimals uint64,
	denom string,
	path string,
	ibcChannelId string,
	ibcCounterpartyChannelId string,
	displayName string,
	displaySymbol string,
	network string,
	address string,
	externalSymbol string,
	transferLimit string,
	permissions []string,
	unitDenom string,
	ibcCounterpartyDenom string,
	ibcCounterpartyChainId string,
) *MsgUpdateEntry {
	return &MsgUpdateEntry{
		Authority:                authority,
		BaseDenom:                baseDenom,
		Decimals:                 decimals,
		Denom:                    denom,
		Path:                     path,
		IbcChannelId:             ibcChannelId,
		IbcCounterpartyChannelId: ibcCounterpartyChannelId,
		DisplayName:              displayName,
		DisplaySymbol:            displaySymbol,
		Network:                  network,
		Address:                  address,
		ExternalSymbol:           externalSymbol,
		TransferLimit:            transferLimit,
		Permissions:              permissions,
		UnitDenom:                unitDenom,
		IbcCounterpartyDenom:     ibcCounterpartyDenom,
		IbcCounterpartyChainId:   ibcCounterpartyChainId,
	}
}

func (msg *MsgUpdateEntry) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if msg.Decimals < 6 || msg.Decimals > 18 {
		return ErrDecimalsInvalid
	}

	if err = sdk.ValidateDenom(msg.BaseDenom); err != nil {
		return ErrInvalidBaseDenom
	}

	if err = sdk.ValidateDenom(msg.Denom); err != nil {
		return fmt.Errorf("invalid denom")
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteEntry{}

func NewMsgDeleteEntry(
	authority string,
	baseDenom string,
) *MsgDeleteEntry {
	return &MsgDeleteEntry{
		Authority: authority,
		BaseDenom: baseDenom,
	}
}

func (msg *MsgDeleteEntry) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if err = sdk.ValidateDenom(msg.BaseDenom); err != nil {
		return ErrInvalidBaseDenom
	}
	return nil
}
