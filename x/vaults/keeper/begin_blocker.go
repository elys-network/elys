package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/vaults/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// Traverse all vaults and deduct management fee from all coins and send it to the vault's manager and protocol revenue address
	// TODO: Add committed tokens
	vaults := k.GetAllVaults(ctx)
	totalBlocksPerYear := k.pk.GetParams(ctx).TotalBlocksPerYear
	protocolAddress := k.masterchef.GetParams(ctx).ProtocolRevenueAddress
	for _, vault := range vaults {
		var protocolCoins sdk.Coins
		coins := k.bk.GetAllBalances(ctx, types.NewVaultAddress(vault.Id))
		var managerCoins sdk.Coins
		for _, coin := range coins {
			coin.Amount = (coin.Amount.ToLegacyDec().Mul(vault.ManagementFee).Quo(math.LegacyNewDecFromInt(math.NewInt(int64(totalBlocksPerYear))))).TruncateInt()

			protocolFeeShare := coin.Amount.ToLegacyDec().Mul(vault.ProtocolFeeShare)
			protocolCoins = protocolCoins.Add(sdk.NewCoin(coin.Denom, protocolFeeShare.TruncateInt()))
			coin.Amount = coin.Amount.Sub(protocolFeeShare.TruncateInt())
			managerCoins = managerCoins.Add(sdk.NewCoin(coin.Denom, coin.Amount))
		}
		// send coins to protocol revenue address and manager address
		err := k.bk.SendCoins(ctx, types.NewVaultAddress(vault.Id), sdk.MustAccAddressFromBech32(vault.Manager), managerCoins)
		if err != nil {
			// log error
			k.Logger(ctx).Error("error sending coins to vault manager", "error", err)
		}
		err = k.bk.SendCoins(ctx, types.NewVaultAddress(vault.Id), sdk.MustAccAddressFromBech32(protocolAddress), protocolCoins)
		if err != nil {
			// log error
			k.Logger(ctx).Error("error sending coins to protocol address", "error", err)
		}
	}

	if k.GetEpochPosition(ctx, k.GetParams(ctx).PerformanceFeeEpochLength) == 0 {
		k.DeductPerformanceFee(ctx)
	}
}

// get position of current block in epoch
func (k Keeper) GetEpochPosition(ctx sdk.Context, epochLength uint64) uint64 {
	if epochLength <= 0 {
		epochLength = 1
	}
	currentHeight := uint64(ctx.BlockHeight())
	return currentHeight % epochLength
}

func (k Keeper) DeductPerformanceFee(ctx sdk.Context) {
	vaults := k.GetAllVaults(ctx)
	totalBlocksPerYear := k.pk.GetParams(ctx).TotalBlocksPerYear
	protocolAddress := k.masterchef.GetParams(ctx).ProtocolRevenueAddress
	for _, vault := range vaults {
		if vault.PerformanceFee.IsPositive() {
			currentValue, err := k.VaultUsdValue(ctx, vault.Id)
			if err != nil {
				k.Logger(ctx).Error("error getting vault value", "error", err)
				continue
			}
			profit := currentValue.Dec().Sub(vault.SumOfDepositsUsdValue).Add(vault.WithdrawalUsdValue)
			if profit.IsPositive() {
				vault.SumOfDepositsUsdValue = vault.SumOfDepositsUsdValue.Add(profit)
				shares := profit.Quo(currentValue.Dec()).Mul(vault.PerformanceFee)

				// TODO: Add committed tokens
				var protocolCoins sdk.Coins
				var managerCoins sdk.Coins
				coins := k.bk.GetAllBalances(ctx, types.NewVaultAddress(vault.Id))
				for _, coin := range coins {
					coin.Amount = (coin.Amount.ToLegacyDec().Mul(shares).Quo(math.LegacyNewDecFromInt(math.NewInt(int64(totalBlocksPerYear))))).TruncateInt()

					protocolFeeShare := coin.Amount.ToLegacyDec().Mul(vault.ProtocolFeeShare)
					protocolCoins = protocolCoins.Add(sdk.NewCoin(coin.Denom, protocolFeeShare.TruncateInt()))
					coin.Amount = coin.Amount.Sub(protocolFeeShare.TruncateInt())
					managerCoins = managerCoins.Add(sdk.NewCoin(coin.Denom, coin.Amount))
				}
				// unwind and send coins to protocol revenue address and manager address
				// send coins to protocol revenue address and manager address
				err := k.bk.SendCoins(ctx, types.NewVaultAddress(vault.Id), sdk.MustAccAddressFromBech32(vault.Manager), managerCoins)
				if err != nil {
					// log error
					k.Logger(ctx).Error("error sending performance fee to vault manager", "error", err)
				}
				err = k.bk.SendCoins(ctx, types.NewVaultAddress(vault.Id), sdk.MustAccAddressFromBech32(protocolAddress), protocolCoins)
				if err != nil {
					// log error
					k.Logger(ctx).Error("error sending performance fee to protocol address", "error", err)
				}
				// TODO: track performance and management fee in state
			}
		}
	}
}

func (k Keeper) EndBlocker(ctx sdk.Context) {
	// Claim rewards for all vaults
	vaults := k.GetAllVaults(ctx)
	for _, vault := range vaults {
		vaultAddress := types.NewVaultAddress(vault.Id)
		vaultRewardCollectorAddress := types.NewVaultRewardCollectorAddress(vault.Id)

		usdcDenom := k.GetBaseCurrencyDenom(ctx)

		beforeBalance := k.commitment.GetAllBalances(ctx, vaultRewardCollectorAddress)
		// TODO: get all pools ids
		err := k.masterchef.ClaimRewards(ctx, vaultAddress, []uint64{1, 2, 3, 4}, vaultRewardCollectorAddress)
		if err != nil {
			k.Logger(ctx).Error("error claiming rewards", "error", err)
		}
		afterBalance := k.commitment.GetAllBalances(ctx, vaultRewardCollectorAddress)
		usdcAmount := afterBalance.AmountOf(usdcDenom).Sub(beforeBalance.AmountOf(usdcDenom))

		// Update reward USDC share
		if usdcAmount.IsPositive() {
			k.UpdateAccPerShare(ctx, vault.Id, usdcDenom, usdcAmount)
		}

		// Update reward EDEN share
		edenAmount := afterBalance.AmountOf(ptypes.Eden).Sub(beforeBalance.AmountOf(ptypes.Eden))
		if edenAmount.IsPositive() {
			k.UpdateAccPerShare(ctx, vault.Id, ptypes.Eden, edenAmount)
		}
	}
}
