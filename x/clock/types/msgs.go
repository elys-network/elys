package types

import (
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Sudo Message called on the contracts
	EndBlockSudoMessage = `{"clock_end_block":{}}`
)

var _ sdk.Msg = &MsgUpdateParams{}

// NewMsgUpdateParams creates new instance of MsgUpdateParams
func NewMsgUpdateParams(
	sender sdk.Address,
	contracts []string,
) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: sender.String(),
		Params:    Params{ContractAddresses: contracts},
	}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrap(err, "invalid authority address")
	}

	return msg.Params.Validate()
}
