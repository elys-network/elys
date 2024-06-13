package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func filterBalancesByDenoms(balances []banktypes.Balance, addressDenomMap map[string][]string) ([]banktypes.Balance, sdk.Coins) {
	newBalances := []banktypes.Balance{}
	var coinsToRemove sdk.Coins

	for _, balance := range balances {
		denomsToKeep, specified := addressDenomMap[balance.Address]
		if !specified {
			// Address is not specified, keep it unchanged
			newBalances = append(newBalances, balance)
			continue
		}

		// Filter the coins for the specified address
		var filteredCoins sdk.Coins
		for _, coin := range balance.Coins {
			if contains(denomsToKeep, coin.Denom) {
				filteredCoins = append(filteredCoins, coin)
			} else {
				coinsToRemove = coinsToRemove.Add(coin)
			}
		}

		if len(filteredCoins) > 0 {
			newBalances = append(newBalances, banktypes.Balance{
				Address: balance.Address,
				Coins:   filteredCoins,
			})
		}
	}

	return newBalances, coinsToRemove
}
