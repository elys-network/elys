package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (c *Commitments) GetUncommittedTokensForDenom(denom string) (*UncommittedTokens, bool) {
	for _, token := range c.UncommittedTokens {
		if token.Denom == denom {
			return token, true
		}
	}
	return &UncommittedTokens{}, false
}

func (c *Commitments) GetCommittedTokensForDenom(denom string) (*CommittedTokens, bool) {
	for _, token := range c.CommittedTokens {
		if token.Denom == denom {
			return token, true
		}
	}
	return &CommittedTokens{}, false
}

func (c *Commitments) GetUncommittedAmountForDenom(denom string) sdk.Int {
	for _, token := range c.UncommittedTokens {
		if token.Denom == denom {
			return token.Amount
		}
	}
	return sdk.NewInt(0)
}

func (c *Commitments) GetCommittedAmountForDenom(denom string) sdk.Int {
	for _, token := range c.CommittedTokens {
		if token.Denom == denom {
			return token.Amount
		}
	}
	return sdk.NewInt(0)
}
