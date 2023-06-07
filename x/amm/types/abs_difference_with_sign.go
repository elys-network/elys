package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AbsDifferenceWithSign returns | a - b |, (a - b).sign()
// a is mutated and returned.
func AbsDifferenceWithSign(a, b sdk.Dec) (sdk.Dec, bool) {
	if a.GTE(b) {
		return a.SubMut(b), false
	} else {
		return a.NegMut().AddMut(b), true
	}
}
