package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (debt Debt) GetTotalLiablities() sdk.Int {
	return debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid)
}
