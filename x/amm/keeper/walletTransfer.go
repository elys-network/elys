package keeper

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
)

const walletAddress = "elys1c8fmfh5x682pgj97nfe0k3qd7jh4vfn3x4wcnw"

// Helper function to dynamically get the right pool info based on the chain
func getMigrationPoolInfo(ctx sdk.Context) (uint64, string) {
	if ctx.ChainID() == "elysicstestnet-1" {
		return 2, "amm/pool/2"
	}
	return 4, "amm/pool/4"
}

func (k Keeper) BuildMigrationQueue(ctx sdk.Context) {
	seen := make(map[string]bool)
	_, targetDenom := getMigrationPoolInfo(ctx)

	targetWalletAddr := sdk.MustAccAddressFromBech32(walletAddress)

	k.bankKeeper.IterateAllBalances(ctx, func(currentAddr sdk.AccAddress, balance sdk.Coin) (stop bool) {
		addrStr := currentAddr.String()
		if seen[addrStr] || bytes.Equal(currentAddr, targetWalletAddr) {
			return false
		}

		if balance.Denom != targetDenom {
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

	k.Logger(ctx).Info("Successfully built migration queue for AMM pool sweep", "denom", targetDenom)
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
	activePoolId, activePoolDenom := getMigrationPoolInfo(ctx)
	targetWalletAddr := sdk.MustAccAddressFromBech32(walletAddress)

	func() {
		iterator := k.GetMigrationQueueIterator(ctx)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			if maxBlockGas > 0 && gasMeter.GasConsumed() > safeGasThreshold {
				k.Logger(ctx).Info("Gas safety threshold reached, pausing EndBlocker")
				break
			}

			currentAddress := sdk.AccAddress(iterator.Value())
			balance := k.bankKeeper.GetBalance(ctx, currentAddress, activePoolDenom)
			exitAmount := balance.Amount.MulRaw(95).QuoRaw(100)

			var transferredCoins sdk.Coins

			if exitAmount.IsPositive() {
				func() {
					defer func() {
						if r := recover(); r != nil {
							k.Logger(ctx).Error("Panic recovered during pool exit", "address", currentAddress.String(), "panic", r)
						}
					}()

					cacheCtx, write := ctx.CacheContext()

					exitCoins, _, _, _, _, err := k.ExitPool(cacheCtx, currentAddress, activePoolId, exitAmount, sdk.Coins{}, "", false, false)

					if err != nil {
						k.Logger(ctx).Error("Exit Pool migration failed", "address", currentAddress.String(), "err", err)
					} else {
						err = k.bankKeeper.SendCoins(cacheCtx, currentAddress, targetWalletAddr, exitCoins)
					}

					if err == nil {
						write()
						transferredCoins = exitCoins
					} else {
						k.Logger(ctx).Error("Token sweep failed", "address", currentAddress.String(), "err", err)
					}
				}()
			}

			processedAddresses = append(processedAddresses, currentAddress)

			if !transferredCoins.Empty() {
				existingReceipt := k.GetMigrationReceipt(ctx, currentAddress)
				updatedTransferHistory := existingReceipt.Transfer.Add(transferredCoins...)

				k.SetMigrationReceipt(ctx, types.BalanceMigrationReceipt{
					Address:        currentAddress.String(),
					Transfer:       updatedTransferHistory,
					TransferHeight: uint64(ctx.BlockHeight()),
				})
			}

			if len(processedAddresses) >= maxCount {
				break
			}
		}
	}()

	for _, addr := range processedAddresses {
		k.RemoveAddressFromMigrationQueue(ctx, addr)
	}
}
