package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
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

	if err = sdk.ValidateDenom(msg.Denom); err != nil {
		return err
	}

	if len(msg.BandTicker) == 0 {
		return fmt.Errorf("band ticker is required")
	}

	if len(msg.ElysTicker) == 0 {
		return fmt.Errorf("elys ticker is required")
	}

	if len(msg.Display) == 0 {
		return fmt.Errorf("display is required")
	}

	if msg.Decimal < 6 || msg.Decimal > 18 {
		return fmt.Errorf("invalid decimal (%d): should be between 6 and 18 (inclusive)", msg.Decimal)
	}
	return nil
}
