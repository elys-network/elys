package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
)

type PerpetualFees struct {
	PerpFees           types.Coins
	SlippageFees       types.Coins
	WeightBreakingFees types.Coins
	TakerFees          types.Coins
}

func NewPerpetualFees(perpFees, slippageFees, weightBreakingFees, takerFees sdk.Coins) PerpetualFees {
	return PerpetualFees{
		PerpFees:           perpFees,
		SlippageFees:       slippageFees,
		WeightBreakingFees: weightBreakingFees,
		TakerFees:          takerFees,
	}
}

func NewPerpetualFeesWithEmptyCoins() PerpetualFees {
	return PerpetualFees{
		PerpFees:           sdk.Coins{},
		SlippageFees:       sdk.Coins{},
		WeightBreakingFees: sdk.Coins{},
		TakerFees:          sdk.Coins{},
	}
}

func (p PerpetualFees) Add(other PerpetualFees) PerpetualFees {
	return NewPerpetualFees(
		p.PerpFees.Add(other.PerpFees...),
		p.SlippageFees.Add(other.SlippageFees...),
		p.WeightBreakingFees.Add(other.WeightBreakingFees...),
		p.TakerFees.Add(other.TakerFees...),
	)
}
