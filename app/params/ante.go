package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Define global variable to be accessible from outside and change these through gov proposal
var MinCommissionRate = sdk.NewDecWithPrec(5, 2) // 5% as a fraction
var MaxVotingPower = sdk.NewDecWithPrec(66, 1)   // 6.6%
var MinSelfDelegation = sdk.NewInt(1)            // 1
