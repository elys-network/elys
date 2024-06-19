package types

import (
	"cosmossdk.io/math"
)

const (
	OneShareExponent = 18
)

var (
	// OneShare represents the amount of subshares in a single pool share.
	OneShare = math.NewIntWithDecimal(1, OneShareExponent)
)
