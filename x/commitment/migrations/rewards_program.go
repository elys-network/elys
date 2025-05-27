package migrations

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/v5/x/commitment/types"
)

var RewardProgram = []types.RewardProgram{
	{
		Address: "elys1234567890",
		Amount:  math.NewInt(1000000000000),
		Claimed: false,
	},
}
