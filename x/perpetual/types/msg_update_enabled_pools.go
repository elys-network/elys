package types

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateEnabledPools{}

func (msg *MsgUpdateEnabledPools) ValidateBasic() error {
	for _, p := range msg.UpdatePoolParams {
		if err := CheckLegacyDecNilAndNegative(p.MtpSafetyFactor, "MtpSafetyFactor"); err != nil {
			return err
		}

		if err := CheckLegacyDecNilAndNegative(p.LeverageMax, "LeverageMax"); err != nil {
			return err
		}

		if p.MtpSafetyFactor.LTE(math.LegacyOneDec()) {
			return fmt.Errorf("MtpSafetyFactor cannot be <= one")
		}

		if p.LeverageMax.LTE(math.LegacyOneDec()) {
			return fmt.Errorf("LeverageMax cannot be <= one")
		}

	}
	return nil
}
