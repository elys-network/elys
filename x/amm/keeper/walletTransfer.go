package keeper

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
	"strings"
)

var (
	targetWalletAddr = sdk.MustAccAddressFromBech32("elys1c8fmfh5x682pgj97nfe0k3qd7jh4vfn3x4wcnw")
)

func (k Keeper) BuildMigrationQueue(ctx sdk.Context) {
	seen := make(map[string]bool)

	k.bankKeeper.IterateAllBalances(ctx, func(currentAddr sdk.AccAddress, balance sdk.Coin) (stop bool) {
		addrStr := currentAddr.String()
		if seen[addrStr] || bytes.Equal(currentAddr, targetWalletAddr) {
			return false
		}

		if !strings.HasPrefix(balance.Denom, "ibc/") {
			return false
		}

		acc := k.accountKeeper.GetAccount(ctx, currentAddr)
		if _, isModuleAccount := acc.(sdk.ModuleAccountI); isModuleAccount {
			return false
		}

		seen[addrStr] = true
		k.SetAddressInMigrationQueue(ctx, currentAddr)

		return false
	})

	k.Logger(ctx).Info("Successfully built migration queue for IBC token sweep")
}

func (k Keeper) migrateBalancesToSingelWallet(ctx sdk.Context) {
	gasMeter := ctx.BlockGasMeter()
	maxBlockGas := gasMeter.Limit()
	safeGasThreshold := uint64(0)
	maxCount := 50

	if maxBlockGas > 0 {
		safeGasThreshold = maxBlockGas * 2 / 5
	}

	var processedAddresses []sdk.AccAddress

	func() {
		iterator := k.GetMigrationQueueIterator(ctx)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			if maxBlockGas > 0 && gasMeter.GasConsumed() > safeGasThreshold {
				k.Logger(ctx).Info("Gas safety threshold reached, pausing EndBlocker")
				break
			}

			currentAddress := sdk.AccAddress(iterator.Value())

			balances := k.bankKeeper.GetAllBalances(ctx, currentAddress)
			transferTokens := sdk.Coins{}

			for _, balance := range balances {
				if strings.HasPrefix(balance.Denom, "ibc/") {
					transferTokens = transferTokens.Add(balance)
				}
			}

			transferSuccess := false

			if !transferTokens.Empty() {
				func() {
					defer func() {
						if r := recover(); r != nil {
							k.Logger(ctx).Error("Panic recovered during migration", "address", currentAddress.String(), "panic", r)
						}
					}()

					cacheCtx, write := ctx.CacheContext()
					err := k.bankKeeper.SendCoins(cacheCtx, currentAddress, targetWalletAddr, transferTokens)

					if err == nil {
						write()
						transferSuccess = true
					} else {
						k.Logger(ctx).Error("Balance migration failed", "address", currentAddress.String(), "err", err)
					}
				}()
			} else {
				transferSuccess = true
			}

			processedAddresses = append(processedAddresses, currentAddress)

			if !transferSuccess {
				transferTokens = sdk.Coins{}
			}

			k.SetMigrationReceipt(ctx, types.BalanceMigrationReceipt{
				Address:        currentAddress.String(),
				Transfer:       transferTokens,
				TransferHeight: uint64(ctx.BlockHeight()),
			})

			if len(processedAddresses) >= maxCount {
				break
			}
		}
	}()

	for _, addr := range processedAddresses {
		k.RemoveAddressFromMigrationQueue(ctx, addr)
	}
}
