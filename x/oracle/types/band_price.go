package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRequestBandPrice = "coin_rates_data"

var (
	_ sdk.Msg = &MsgRequestBandPrice{}

	// BandPriceResultStoreKeyPrefix is a prefix for storing result
	BandPriceResultStoreKeyPrefix = "coin_rates_result"

	// LastBandRequestIdKey is the key for the last request id
	LastBandRequestIdKey = "coin_rates_last_id"

	// BandPriceClientIDKey is query request identifier
	BandPriceClientIDKey = "coin_rates_id"

	// PrefixKeyBandRequest is the prefix for band requests
	PrefixKeyBandRequest = "band_request_"
)

// NewMsgRequestBandPrice creates a new BandPrice message
func NewMsgRequestBandPrice(
	creator string,
	oracleScriptID OracleScriptID,
	sourceChannel string,
	calldata *BandPriceCallData,
	askCount uint64,
	minCount uint64,
	feeLimit sdk.Coins,
	prepareGas uint64,
	executeGas uint64,
) *MsgRequestBandPrice {
	return &MsgRequestBandPrice{
		ClientID:       BandPriceClientIDKey,
		Creator:        creator,
		OracleScriptID: uint64(oracleScriptID),
		SourceChannel:  sourceChannel,
		Calldata:       calldata,
		AskCount:       askCount,
		MinCount:       minCount,
		FeeLimit:       feeLimit,
		PrepareGas:     prepareGas,
		ExecuteGas:     executeGas,
	}
}

// Route returns the message route
func (m *MsgRequestBandPrice) Route() string {
	return RouterKey
}

// Type returns the message type
func (m *MsgRequestBandPrice) Type() string {
	return TypeMsgRequestBandPrice
}

// GetSigners returns the message signers
func (m *MsgRequestBandPrice) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the signed bytes from the message
func (m *MsgRequestBandPrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic check the basic message validation
func (m *MsgRequestBandPrice) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if m.SourceChannel == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid source channel")
	}
	return nil
}

// BandPriceResultStoreKey is a function to generate key for each result in store
func BandPriceResultStoreKey(requestID OracleRequestID) []byte {
	return append(KeyPrefix(BandPriceResultStoreKeyPrefix), int64ToBytes(int64(requestID))...)
}

func BandRequestStoreKey(requestID OracleRequestID) []byte {
	return append(KeyPrefix(PrefixKeyBandRequest), int64ToBytes(int64(requestID))...)
}
