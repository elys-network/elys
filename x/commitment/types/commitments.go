package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
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
			if unlockTime != 0 {
				c.CommittedTokens[i].Lockups = append(token.Lockups, Lockup{
					Amount:          amount,
					UnlockTimestamp: unlockTime,
				})
			}
			return
		}
	}

	committedToken := &CommittedTokens{
		Denom:   denom,
		Amount:  amount,
		Lockups: []Lockup{},
	}
	if unlockTime != 0 {
		committedToken.Lockups = append(committedToken.Lockups, Lockup{
			Amount:          amount,
			UnlockTimestamp: unlockTime,
		})
	}
	c.CommittedTokens = append(c.CommittedTokens, committedToken)
}

func (c Commitments) CommittedTokensLocked(ctx sdk.Context) (sdk.Coins, sdk.Coins) {
	totalLocked := sdk.Coins{}
	totalCommitted := sdk.Coins{}
	for _, token := range c.CommittedTokens {
		lockedAmount := sdk.ZeroInt()
		for _, lockup := range token.Lockups {
			if lockup.UnlockTimestamp > uint64(ctx.BlockTime().Unix()) {
				lockedAmount = lockedAmount.Add(lockup.Amount)
			}
		}
		totalLocked = totalLocked.Add(sdk.NewCoin(token.Denom, lockedAmount))
		totalCommitted = totalCommitted.Add(sdk.NewCoin(token.Denom, token.Amount))
	}
	return totalLocked, totalCommitted
}

func (c *Commitments) DeductFromCommitted(denom string, amount math.Int, currTime uint64) error {
	for i, token := range c.CommittedTokens {
		if token.Denom == denom {
			c.CommittedTokens[i].Amount = token.Amount.Sub(amount)
			if c.CommittedTokens[i].Amount.IsNegative() {
				return ErrInsufficientCommittedTokens
			}

			newLockups := []Lockup{}
			lockedAmount := sdk.ZeroInt()
			for _, lockup := range token.Lockups {
				if lockup.UnlockTimestamp > currTime {
					newLockups = append(newLockups, lockup)
					lockedAmount = lockedAmount.Add(lockup.Amount)
				}
			}
			c.CommittedTokens[i].Lockups = newLockups
			if lockedAmount.GT(c.CommittedTokens[i].Amount) {
				return errors.Wrapf(ErrInsufficientWithdrawableTokens, "amount: %s denom: %s", amount, denom)
			}
			if c.CommittedTokens[i].Amount.IsZero() {
				c.CommittedTokens = append(c.CommittedTokens[:i], c.CommittedTokens[i+1:]...)
			}
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

func (c Commitments) GetCreatorAccount() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(c.Creator)
}

func (vesting *VestingTokens) VestedSoFar(ctx sdk.Context) math.Int {
	totalBlocks := ctx.BlockHeight() - vesting.StartBlock
	if totalBlocks > vesting.NumBlocks {
		totalBlocks = vesting.NumBlocks
	}
	return vesting.TotalAmount.Mul(sdk.NewInt(totalBlocks)).Quo(sdk.NewInt(vesting.NumBlocks))
}
