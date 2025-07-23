package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v6/utils"
)

var _ sdk.Msg = &MsgClose{}

func NewMsgClose(creator string, id uint64, amount math.Int, poolId uint64, closingRatio math.LegacyDec) *MsgClose {
	return &MsgClose{
		Creator:      creator,
		Id:           id,
		Amount:       amount,
		PoolId:       poolId,
		ClosingRatio: closingRatio,
	}
}

func (msg *MsgClose) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Amount.IsNil() || msg.Amount.IsNegative() {
		return ErrInvalidAmount
	}

	if err = utils.CheckLegacyDecNilAndNegative(msg.ClosingRatio, "ClosingRatio"); err != nil {
		return err
	}

	if msg.Amount.IsZero() && msg.ClosingRatio.IsZero() {
		return errors.New("closing ratio and amount both cannot be zero")
	}

	if msg.PoolId == 0 {
		return errors.New("invalid pool id")
	}

	return nil
}
