package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateRewardsDataLifetime{}

func NewMsgUpdateRewardsDataLifetime(creator string, rewardsDataLifetime string) *MsgUpdateRewardsDataLifetime {
	return &MsgUpdateRewardsDataLifetime{
		Creator:             creator,
		RewardsDataLifetime: rewardsDataLifetime,
	}
}

func (msg *MsgUpdateRewardsDataLifetime) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	res, ok := math.NewIntFromString(msg.RewardsDataLifetime)
	if !ok || !res.IsPositive() {
		return errorsmod.Wrapf(ErrInvalidRewardsDataLifecycle, "invalid data in rewards_data_lifecycle")
	}
	return nil
}
