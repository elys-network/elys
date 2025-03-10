package types

import (
	fmt "fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdatePool{}

func NewMsgUpdatePool(signer string, interestRateMax math.LegacyDec, interestRateMin math.LegacyDec, interestRateIncrease math.LegacyDec,
	interestRateDecrease math.LegacyDec, healthFactor math.LegacyDec, maxLeverageRatio math.LegacyDec, poolId uint64) *MsgUpdatePool {
	return &MsgUpdatePool{
		Authority:            signer,
		InterestRateMax:      interestRateMax,
		InterestRateMin:      interestRateMin,
		InterestRateIncrease: interestRateIncrease,
		InterestRateDecrease: interestRateDecrease,
		HealthGainFactor:     healthFactor,
		MaxLeverageRatio:     maxLeverageRatio,
		PoolId:               poolId,
	}
}

func (msg *MsgUpdatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := CheckLegacyDecNilAndNegative(msg.InterestRateMax, "InterestRateMax"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(msg.InterestRateMin, "InterestRateMin"); err != nil {
		return err
	}
	if msg.InterestRateMax.LT(msg.InterestRateMin) {
		return fmt.Errorf("InterestRateMax must be greater than InterestRateMin")
	}
	if err := CheckLegacyDecNilAndNegative(msg.InterestRateIncrease, "InterestRateIncrease"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(msg.InterestRateDecrease, "InterestRateDecrease"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(msg.HealthGainFactor, "HealthGainFactor"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(msg.MaxLeverageRatio, "MaxLeverageRatio"); err != nil {
		return err
	}

	return nil
}
