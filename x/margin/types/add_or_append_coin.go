package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Utility function to add a new coin or aggregate the amount if coin with same denom already exists
func AddOrAppendCoin(coins []sdk.Coin, coin sdk.Coin) []sdk.Coin {
	for i, c := range coins {
		if c.Denom == coin.Denom {
			coins[i].Amount = c.Amount.Add(coin.Amount)
			return coins
		}
	}
	return append(coins, coin)
}
