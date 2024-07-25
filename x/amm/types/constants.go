package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	OneShareExponent = 18

	BalancerGasFeeForSwap = 10_000
)

var (
	// OneShare represents the amount of subshares in a single pool share.
	OneShare = math.NewIntWithDecimal(1, OneShareExponent)

	// InitPoolSharesSupply is the amount of new shares to initialize a pool with.
	InitPoolSharesSupply = OneShare.MulRaw(100)

	// GuaranteedWeightPrecision Scaling factor for every weight. The pool weight is:
	// weight_in_MsgCreateBalancerPool * GuaranteedWeightPrecision
	//
	// This is done so that smooth weight changes have enough precision to actually be smooth.
	GuaranteedWeightPrecision int64 = 1 << 30

	oneHalf    = sdk.MustNewDecFromStr("0.5")
	twoDec     = sdk.MustNewDecFromStr("2")
	ln2        = sdk.MustNewDecFromStr("0.693147180559945309")
	inverseLn2 = sdk.MustNewDecFromStr("1.442695040888963407")
	exp        = sdk.MustNewDecFromStr("2.718281828459045235")

	// PowPrecision Don't EVER change after initializing
	// TODO: Analyze choice here.
	powPrecision = sdk.MustNewDecFromStr("0.00000001")
)
