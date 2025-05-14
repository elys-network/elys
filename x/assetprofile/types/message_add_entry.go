package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddEntry{}

func NewMsgAddEntry(creator string, baseDenom string, decimals uint64, denom string, path string, ibcChannelId string, ibcCounterpartyChannelId string, displayName string, displaySymbol string, network string, address string, externalSymbol string, transferLimit string, permissions []string, unitDenom string, ibcCounterpartyDenom string, ibcCounterpartyChainId string, commitEnabled bool, withdrawEnabled bool) *MsgAddEntry {
	return &MsgAddEntry{
		Creator:                  creator,
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
		CommitEnabled:            commitEnabled,
		WithdrawEnabled:          withdrawEnabled,
	}
}

func (msg *MsgAddEntry) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Decimals < 6 || msg.Decimals > 18 {
		return ErrDecimalsInvalid
	}

	// eureka assets might have different base denoms
	//if err = sdk.ValidateDenom(msg.BaseDenom); err != nil {
	//	return ErrInvalidBaseDenom
	//}
	if msg.BaseDenom == "" {
		return errorsmod.Wrapf(ErrInvalidBaseDenom, "base denom cannot be empty")
	}
	if len(msg.BaseDenom) == 1 || len(msg.BaseDenom) > 128 {
		return errorsmod.Wrapf(ErrInvalidBaseDenom, "base denom cannot be longer than 128 characters or one character")
	}

	if err = sdk.ValidateDenom(msg.Denom); err != nil {
		return err
	}

	return nil
}
