package types

import (
	"cosmossdk.io/math"
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

	// Scaling factor for every weight. The pool weight is:
	// weight_in_MsgCreateBalancerPool * GuaranteedWeightPrecision
	//
	// This is done so that smooth weight changes have enough precision to actually be smooth.
	GuaranteedWeightPrecision int64 = 1 << 30
)
