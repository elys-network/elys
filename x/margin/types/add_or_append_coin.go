package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Utility function to add a new coin or aggregate the amount if coin with same denom already exists
func AddOrAppendCoin(coins sdk.Coins, coin sdk.Coin) sdk.Coins {
	return coins.Add(coin)
}
