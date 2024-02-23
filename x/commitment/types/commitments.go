package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (c Commitments) IsEmpty() bool {
	if len(c.CommittedTokens) > 0 {
		return false
	}
	if len(c.RewardsUnclaimed) > 0 {
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

func (c *Commitments) GetRewardUnclaimedForDenom(denom string) math.Int {
	return c.RewardsUnclaimed.AmountOf(denom)
}

// Sub bucket rewards query - Elys
func (c *Commitments) GetElysSubBucketRewardUnclaimedForDenom(denom string) math.Int {
	return c.RewardsByElysUnclaimed.AmountOf(denom)
}

// Sub bucket rewards query - Eden
func (c *Commitments) GetEdenSubBucketRewardUnclaimedForDenom(denom string) math.Int {
	return c.RewardsByEdenUnclaimed.AmountOf(denom)
}

// Sub bucket rewards query - EdenB
func (c *Commitments) GetEdenBSubBucketRewardUnclaimedForDenom(denom string) math.Int {
	return c.RewardsByEdenbUnclaimed.AmountOf(denom)
}

// Sub bucket rewards query - Usdc
func (c *Commitments) GetUsdcSubBucketRewardUnclaimedForDenom(denom string) math.Int {
	return c.RewardsByUsdcUnclaimed.AmountOf(denom)
}

func (c *Commitments) AddRewardsUnclaimed(amount sdk.Coin) {
	c.RewardsUnclaimed = c.RewardsUnclaimed.Add(amount)
}

func (c *Commitments) AddSubBucketRewardsByElysUnclaimed(amount sdk.Coins) {
	c.RewardsByElysUnclaimed = c.RewardsByElysUnclaimed.Add(amount...)
}

func (c *Commitments) AddSubBucketRewardsByEdenUnclaimed(amount sdk.Coins) {
	c.RewardsByEdenUnclaimed = c.RewardsByEdenUnclaimed.Add(amount...)
}

func (c *Commitments) AddSubBucketRewardsByEdenBUnclaimed(amount sdk.Coins) {
	c.RewardsByEdenbUnclaimed = c.RewardsByEdenbUnclaimed.Add(amount...)
}

func (c *Commitments) AddSubBucketRewardsByUsdcUnclaimed(amount sdk.Coins) {
	c.RewardsByUsdcUnclaimed = c.RewardsByUsdcUnclaimed.Add(amount...)
}

func (c *Commitments) SubRewardsUnclaimed(amount sdk.Coin) error {
	if c.RewardsUnclaimed.AmountOf(amount.Denom).LT(amount.Amount) {
		return ErrInsufficientRewardsUnclaimed
	}
	c.RewardsUnclaimed = c.RewardsUnclaimed.Sub(amount)

	// First deduct from elys staking bucket
	allDeducted, remainedAmount := c.DeductSubBucketRewardsByElysUnclaimed(amount)
	if allDeducted {
		return nil
	}

	// If there is still remaining amount, deduct from eden commited bucket
	allDeducted, remainedAmount = c.DeductSubBucketRewardsByEdenUnclaimed(remainedAmount)
	if allDeducted {
		return nil
	}

	// If there is still remaining amount, deduct from edenb commited bucket
	allDeducted, remainedAmount = c.DeductSubBucketRewardsByEdenBUnclaimed(remainedAmount)
	if allDeducted {
		return nil
	}

	// If there is still remaining amount, deduct from usdc deposit bucket
	c.DeductSubBucketRewardsByUsdcUnclaimed(remainedAmount)

	return nil
}

// This is the function used when we withdraw the reward by program
func (c *Commitments) SubRewardsUnclaimedForElysStaking(amount sdk.Coin) error {
	if c.RewardsUnclaimed.AmountOf(amount.Denom).LT(amount.Amount) {
		return ErrInsufficientRewardsUnclaimed
	}
	c.RewardsUnclaimed = c.RewardsUnclaimed.Sub(amount)

	c.DeductSubBucketRewardsByElysUnclaimed(amount)
	return nil
}

// This is the function used when we withdraw the reward by program
func (c *Commitments) SubRewardsUnclaimedForEdenCommitted(amount sdk.Coin) error {
	if c.RewardsUnclaimed.AmountOf(amount.Denom).LT(amount.Amount) {
		return ErrInsufficientRewardsUnclaimed
	}
	c.RewardsUnclaimed = c.RewardsUnclaimed.Sub(amount)

	c.DeductSubBucketRewardsByEdenUnclaimed(amount)
	return nil
}

// This is the function used when we withdraw the reward by program
func (c *Commitments) SubRewardsUnclaimedForEdenBCommitted(amount sdk.Coin) error {
	if c.RewardsUnclaimed.AmountOf(amount.Denom).LT(amount.Amount) {
		return ErrInsufficientRewardsUnclaimed
	}
	c.RewardsUnclaimed = c.RewardsUnclaimed.Sub(amount)

	c.DeductSubBucketRewardsByEdenBUnclaimed(amount)
	return nil
}

// This is the function used when we withdraw the reward by program
func (c *Commitments) SubRewardsUnclaimedForUSDCDeposit(amount sdk.Coin) error {
	if c.RewardsUnclaimed.AmountOf(amount.Denom).LT(amount.Amount) {
		return ErrInsufficientRewardsUnclaimed
	}
	c.RewardsUnclaimed = c.RewardsUnclaimed.Sub(amount)

	c.DeductSubBucketRewardsByUsdcUnclaimed(amount)
	return nil
}

func (c *Commitments) DeductSubBucketRewardsByElysUnclaimed(amount sdk.Coin) (bool, sdk.Coin) {
	amountToDeduct := amount
	availableToDeduct := c.RewardsByElysUnclaimed.AmountOf(amount.Denom)
	if availableToDeduct.LT(amountToDeduct.Amount) {
		amountToDeduct = sdk.NewCoin(amount.Denom, availableToDeduct)
	}

	c.RewardsByElysUnclaimed = c.RewardsByElysUnclaimed.Sub(amountToDeduct)
	remainedAmount := amount.Sub(amountToDeduct)
	if remainedAmount.Amount.LTE(sdk.ZeroInt()) {
		return true, remainedAmount
	}

	return false, remainedAmount
}

func (c *Commitments) DeductSubBucketRewardsByEdenUnclaimed(amount sdk.Coin) (bool, sdk.Coin) {
	amountToDeduct := amount
	availableToDeduct := c.RewardsByEdenUnclaimed.AmountOf(amount.Denom)
	if availableToDeduct.LT(amountToDeduct.Amount) {
		amountToDeduct = sdk.NewCoin(amount.Denom, availableToDeduct)
	}

	c.RewardsByEdenUnclaimed = c.RewardsByEdenUnclaimed.Sub(amountToDeduct)
	remainedAmount := amount.Sub(amountToDeduct)
	if remainedAmount.Amount.LTE(sdk.ZeroInt()) {
		return true, remainedAmount
	}

	return false, remainedAmount
}

func (c *Commitments) DeductSubBucketRewardsByEdenBUnclaimed(amount sdk.Coin) (bool, sdk.Coin) {
	amountToDeduct := amount
	availableToDeduct := c.RewardsByEdenbUnclaimed.AmountOf(amount.Denom)
	if availableToDeduct.LT(amountToDeduct.Amount) {
		amountToDeduct = sdk.NewCoin(amount.Denom, availableToDeduct)
	}

	c.RewardsByEdenbUnclaimed = c.RewardsByEdenbUnclaimed.Sub(amountToDeduct)
	remainedAmount := amount.Sub(amountToDeduct)
	if remainedAmount.Amount.LTE(sdk.ZeroInt()) {
		return true, remainedAmount
	}

	return false, remainedAmount
}

func (c *Commitments) DeductSubBucketRewardsByUsdcUnclaimed(amount sdk.Coin) (bool, sdk.Coin) {
	amountToDeduct := amount
	availableToDeduct := c.RewardsByUsdcUnclaimed.AmountOf(amount.Denom)
	if availableToDeduct.LT(amountToDeduct.Amount) {
		amountToDeduct = sdk.NewCoin(amount.Denom, availableToDeduct)
	}

	c.RewardsByUsdcUnclaimed = c.RewardsByUsdcUnclaimed.Sub(amountToDeduct)
	remainedAmount := amount.Sub(amountToDeduct)
	if remainedAmount.Amount.LTE(sdk.ZeroInt()) {
		return true, remainedAmount
	}

	return false, remainedAmount
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
