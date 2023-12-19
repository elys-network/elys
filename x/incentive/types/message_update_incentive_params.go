package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateIncentiveParams = "update_incentive_params"

var _ sdk.Msg = &MsgUpdateIncentiveParams{}

func NewMsgUpdateIncentiveParams(creator string, communityTax sdk.Dec, withdrawAddrEnabled bool, rewardPortionForLps sdk.Dec, rewardPortionForStakers sdk.Dec, elysStakeTrackingRate int64, maxEdenRewardAprStakers sdk.Dec, maxEdenRewardParLps sdk.Dec, distributionEpochForStakers int64, distributionEpochForLps int64) *MsgUpdateIncentiveParams {
	return &MsgUpdateIncentiveParams{
		Authority:                   creator,
		CommunityTax:                communityTax,
		WithdrawAddrEnabled:         withdrawAddrEnabled,
		RewardPortionForLps:         rewardPortionForLps,
		RewardPortionForStakers:     rewardPortionForStakers,
		ElysStakeTrackingRate:       elysStakeTrackingRate,
		MaxEdenRewardAprStakers:     maxEdenRewardAprStakers,
		MaxEdenRewardAprLps:         maxEdenRewardParLps,
		DistributionEpochForStakers: distributionEpochForStakers,
		DistributionEpochForLps:     distributionEpochForLps,
	}
}

func (msg *MsgUpdateIncentiveParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdateIncentiveParams) Type() string {
	return TypeMsgUpdateIncentiveParams
}

func (msg *MsgUpdateIncentiveParams) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateIncentiveParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateIncentiveParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.RewardPortionForLps.GT(sdk.NewDec(1)) {
		return errorsmod.Wrapf(sdkerrors.ErrNotSupported, "invalid rewards portion for LPs (%s)", errors.New("Invalid LP portion"))
	}
	if msg.RewardPortionForStakers.GT(sdk.NewDec(1)) {
		return errorsmod.Wrapf(sdkerrors.ErrNotSupported, "invalid rewards portion for Stakers (%s)", errors.New("Invalid Staker portion"))
	}
	if msg.RewardPortionForLps.Add(msg.RewardPortionForStakers).GT(sdk.NewDec(1)) {
		return errorsmod.Wrapf(sdkerrors.ErrNotSupported, "invalid rewards portion for Stakers and LPs (%s)", errors.New("Invalid Staker and LP portion"))
	}
	if msg.MaxEdenRewardAprStakers.LT(sdk.ZeroDec()) {
		return errorsmod.Wrapf(sdkerrors.ErrNotSupported, "invalid max eden rewards apr for stakers (%s)", errors.New("Invalid Rewards APR"))
	}
	if msg.MaxEdenRewardAprLps.LT(sdk.ZeroDec()) {
		return errorsmod.Wrapf(sdkerrors.ErrNotSupported, "invalid max eden rewards apr for stakers (%s)", errors.New("Invalid Rewards APR"))
	}
	if msg.DistributionEpochForStakers < 1 {
		return errorsmod.Wrapf(sdkerrors.ErrNotSupported, "invalid distribution epoch (%s)", errors.New("Invalid epoch"))
	}
	if msg.DistributionEpochForLps < 1 {
		return errorsmod.Wrapf(sdkerrors.ErrNotSupported, "invalid distribution epoch (%s)", errors.New("Invalid epoch"))
	}
	if msg.ElysStakeTrackingRate < 1 {
		return errorsmod.Wrapf(sdkerrors.ErrNotSupported, "invalid elys staked tracking epoch (%s)", errors.New("Invalid elys staked tracking epoch"))
	}
	return nil
}
