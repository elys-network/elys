package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCoinRatesData = "coin_rates_data"

var (
	_ sdk.Msg = &MsgCoinRatesData{}

	// CoinRatesResultStoreKeyPrefix is a prefix for storing result
	CoinRatesResultStoreKeyPrefix = "coin_rates_result"

	// LastCoinRatesIDKey is the key for the last request id
	LastCoinRatesIDKey = "coin_rates_last_id"

	// CoinRatesClientIDKey is query request identifier
	CoinRatesClientIDKey = "coin_rates_id"
)

// NewMsgCoinRatesData creates a new CoinRates message
func NewMsgCoinRatesData(
	creator string,
	oracleScriptID OracleScriptID,
	sourceChannel string,
	calldata *CoinRatesCallData,
	askCount uint64,
	minCount uint64,
	feeLimit sdk.Coins,
	prepareGas uint64,
	executeGas uint64,
) *MsgCoinRatesData {
	return &MsgCoinRatesData{
		ClientID:       CoinRatesClientIDKey,
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
func (m *MsgCoinRatesData) Route() string {
	return RouterKey
}

// Type returns the message type
func (m *MsgCoinRatesData) Type() string {
	return TypeMsgCoinRatesData
}

// GetSigners returns the message signers
func (m *MsgCoinRatesData) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the signed bytes from the message
func (m *MsgCoinRatesData) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic check the basic message validation
func (m *MsgCoinRatesData) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if m.SourceChannel == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid source channel")
	}
	return nil
}

// CoinRatesResultStoreKey is a function to generate key for each result in store
func CoinRatesResultStoreKey(requestID OracleRequestID) []byte {
	return append(KeyPrefix(CoinRatesResultStoreKeyPrefix), int64ToBytes(int64(requestID))...)
}
