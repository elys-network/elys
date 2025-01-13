package types

import (
	fmt "fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddPool{}

func NewMsgAddPool(signer string, depositDenom string, interestRate math.LegacyDec, interestRateMax math.LegacyDec, interestRateMin math.LegacyDec, interestRateIncrease math.LegacyDec,
	interestRateDecrease math.LegacyDec, healthFactor math.LegacyDec, maxLeverageRatio math.LegacyDec, poolId uint64) *MsgAddPool {
	return &MsgAddPool{
		Authority:            signer,
		DepositDenom:         depositDenom,
		InterestRate:         interestRate,
		InterestRateMax:      interestRateMax,
		InterestRateMin:      interestRateMin,
		InterestRateIncrease: interestRateIncrease,
		InterestRateDecrease: interestRateDecrease,
		HealthGainFactor:     healthFactor,
		MaxLeverageRatio:     maxLeverageRatio,
	}
}

func (msg *MsgAddPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := sdk.ValidateDenom(msg.DepositDenom); err != nil {
		return err
	}

	if err := CheckLegacyDecNilAndNegative(msg.InterestRate, "InterestRate"); err != nil {
		return err
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
