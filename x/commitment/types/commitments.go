package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (c *Commitments) GetCommittedAmountForDenom(denom string) sdk.Int {
	for _, token := range c.CommittedTokens {
		if token.Denom == denom {
			return token.Amount
		}
	}
	return sdk.NewInt(0)
}

func (c *Commitments) AddCommittedTokens(denom string, amount sdk.Int, unlockTime uint64) {
	for i, token := range c.CommittedTokens {
		if token.Denom == denom {
			c.CommittedTokens[i].Amount = token.Amount.Add(amount)
			c.CommittedTokens[i].Lockups = append(token.Lockups, Lockup{
				Amount:          amount,
				UnlockTimestamp: unlockTime,
			})

			li := int(0)
			for _, lockup := range token.Lockups {
				if lockup.UnlockTimestamp < unlockTime {
					li++
					continue
				} else if lockup.UnlockTimestamp == unlockTime {
					c.CommittedTokens[i].Lockups[li].Amount = lockup.Amount.Add(amount)
					return
				} else {
					break
				}
			}
			c.CommittedTokens[i].Lockups = append(token.Lockups[:li+1], token.Lockups[li:]...)
			c.CommittedTokens[i].Lockups[li] = Lockup{
				Amount:          amount,
				UnlockTimestamp: unlockTime,
			}
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

func (c *Commitments) DeductFromCommitted(denom string, amount sdk.Int, currTime uint64) error {
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

func (c *Commitments) GetRewardUnclaimedForDenom(denom string) sdk.Int {
	for _, token := range c.RewardsUnclaimed {
		if token.Denom == denom {
			return token.Amount
		}
	}
	return sdk.ZeroInt()
}

func (c *Commitments) AddRewardsUnclaimed(amount sdk.Coin) {
	c.RewardsUnclaimed = c.RewardsUnclaimed.Add(amount)
}

func (c *Commitments) SubRewardsUnclaimed(amount sdk.Coin) error {
	if c.RewardsUnclaimed.AmountOf(amount.Denom).LT(amount.Amount) {
		return ErrInsufficientRewardsUnclaimed
	}
	c.RewardsUnclaimed = c.RewardsUnclaimed.Sub(amount)
	return nil
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
