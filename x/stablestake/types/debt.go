package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (debt Debt) GetTotalLiablities() sdkmath.Int {
	return debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid)
}

func (debt Debt) GetBigDecTotalLiablities() osmomath.BigDec {
	return osmomath.BigDecFromSDKInt(debt.GetTotalLiablities())
}

func (debt Debt) GetOwnerAccount() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(debt.Address)
}

func (debt Debt) GetBigDecBorrowed() osmomath.BigDec {
	return osmomath.BigDecFromSDKInt(debt.Borrowed)
}
