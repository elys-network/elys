package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (debt Debt) GetTotalLiablities() sdk.Int {
	return debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid)
}

func (debt Debt) GetOwnerAccount() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(debt.Address)
}
