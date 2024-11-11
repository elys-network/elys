package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (debt Debt) GetTotalLiablities() sdkmath.Int {
	return debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid)
}

func (debt Debt) GetOwnerAccount() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(debt.Address)
}
