package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (c Commitments) IsEmpty() bool {
	if len(c.CommittedTokens) > 0 {
		return false
	}
	if len(c.Claimed) > 0 {
		return false
	}
	if len(c.VestingTokens) > 0 {
		return false
	}
	return true
}

func (c *Commitments) GetCommittedAmountForDenom(denom string) math.Int {
	for _, token := range c.CommittedTokens {
		if token.Denom == denom {
			return token.Amount
		}
	}
	return sdk.NewInt(0)
}

func (c *Commitments) GetCommittedLockUpsForDenom(denom string) []Lockup {
	for _, token := range c.CommittedTokens {
		if token.Denom == denom {
			return token.Lockups
		}
	}
	return nil
}

func (c *Commitments) AddCommittedTokens(denom string, amount math.Int, unlockTime uint64) {
	for i, token := range c.CommittedTokens {
		if token.Denom == denom {
			c.CommittedTokens[i].Amount = token.Amount.Add(amount)
			c.CommittedTokens[i].Lockups = append(token.Lockups, Lockup{
				Amount:          amount,
				UnlockTimestamp: unlockTime,
			})
			return
		}
	}
	c.CommittedTokens = append(c.CommittedTokens, &CommittedTokens{
		Denom:  denom,
		Amount: amount,
		Lockups: []Lockup{
			{
				Amount:          amount,
				UnlockTimestamp: unlockTime,
			},
		},
	})
}

func (c *Commitments) DeductFromCommitted(denom string, amount math.Int, currTime uint64) error {
	for i, token := range c.CommittedTokens {
		if token.Denom == denom {
			c.CommittedTokens[i].Amount = token.Amount.Sub(amount)
			if c.CommittedTokens[i].Amount.IsNegative() {
				return ErrInsufficientCommittedTokens
			}

			withdrawnTokens := sdk.ZeroInt()
			newLockups := []Lockup{}

			for _, lockup := range token.Lockups {
				if withdrawnTokens.LT(amount) {
					if lockup.UnlockTimestamp <= currTime {
						withdrawAmount := lockup.Amount
						if withdrawAmount.GT(amount.Sub(withdrawnTokens)) {
							withdrawAmount = amount.Sub(withdrawnTokens)
							newLockups = append(newLockups, Lockup{
								Amount:          lockup.Amount.Sub(withdrawAmount),
								UnlockTimestamp: lockup.UnlockTimestamp,
							})
						}
						withdrawnTokens = withdrawnTokens.Add(withdrawAmount)
					} else {
						return ErrInsufficientWithdrawableTokens
					}
				} else {
					newLockups = append(newLockups, lockup)
				}
			}
			c.CommittedTokens[i].Lockups = newLockups
			return nil
		}
	}
	return ErrInsufficientCommittedTokens
}

func (c *Commitments) GetClaimedForDenom(denom string) math.Int {
	for _, token := range c.Claimed {
		if token.Denom == denom {
			return token.Amount
		}
	}
	return sdk.ZeroInt()
}

func (c *Commitments) AddClaimed(amount sdk.Coin) {
	c.Claimed = c.Claimed.Add(amount)
}

func (c *Commitments) SubClaimed(amount sdk.Coin) error {
	if c.Claimed.AmountOf(amount.Denom).LT(amount.Amount) {
		return ErrInsufficientClaimed
	}
	c.Claimed = c.Claimed.Sub(amount)
	return nil
}

func (vesting *VestingTokens) VestedSoFar(ctx sdk.Context) math.Int {
	totalBlocks := ctx.BlockHeight() - vesting.StartBlock
	if totalBlocks > vesting.NumBlocks {
		totalBlocks = vesting.NumBlocks
	}
	return vesting.TotalAmount.Mul(sdk.NewInt(totalBlocks)).Quo(sdk.NewInt(vesting.NumBlocks))
}
