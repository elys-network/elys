package keeper

import (
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/vaults/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// Traverse all vaults and deduct management fee from all coins and send it to the vault's manager and protocol revenue address

	vaults := k.GetAllVaults(ctx)
	totalBlocksPerYear := k.pk.GetParams(ctx).TotalBlocksPerYear
	protocolAddress := k.masterchef.GetParams(ctx).ProtocolRevenueAddress
	for _, vault := range vaults {
		k.distributeVaultFees(ctx, vault, vault.ManagementFee, totalBlocksPerYear, protocolAddress)
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
				k.distributeVaultFees(ctx, vault, shares, totalBlocksPerYear, protocolAddress)
				// TODO: track performance and management fee in state
			}
		}
	}
}

func (k Keeper) EndBlocker(ctx sdk.Context) {
	if k.GetEpochPosition(ctx, k.GetParams(ctx).PerformanceFeeEpochLength) == 0 {
		// Claim rewards for all vaults
		vaults := k.GetAllVaults(ctx)
		for _, vault := range vaults {
			vaultAddress := types.NewVaultAddress(vault.Id)
			vaultRewardCollectorAddress := types.NewVaultRewardCollectorAddress(vault.Id)

			usdcDenom := k.GetBaseCurrencyDenom(ctx)

			beforeBalance := k.commitment.GetAllBalances(ctx, vaultRewardCollectorAddress)
			poolIds := k.GetAllPoolIds(ctx, vaultAddress)
			err := k.masterchef.ClaimRewards(ctx, vaultAddress, poolIds, vaultRewardCollectorAddress)
			if err != nil {
				k.Logger(ctx).Error("error claiming rewards", "error", err)
			}
			afterBalance := k.commitment.GetAllBalances(ctx, vaultRewardCollectorAddress)
			usdcAmount := afterBalance.AmountOf(usdcDenom).Sub(beforeBalance.AmountOf(usdcDenom))

			// Send usdc to vault address
			if usdcAmount.IsPositive() {
				err = k.bk.SendCoins(ctx, vaultRewardCollectorAddress, vaultAddress, sdk.NewCoins(sdk.NewCoin(usdcDenom, usdcAmount)))
				if err != nil {
					k.Logger(ctx).Error("error sending usdc to vault address", "error", err)
				}
			}

			// Update reward EDEN share
			edenAmount := afterBalance.AmountOf(ptypes.Eden).Sub(beforeBalance.AmountOf(ptypes.Eden))
			if edenAmount.IsPositive() {
				k.UpdateAccPerShare(ctx, vault.Id, ptypes.Eden, edenAmount)
			}

			k.AddPoolRewardsAccum(
				ctx,
				vault.Id,
				uint64(ctx.BlockTime().Unix()),
				ctx.BlockHeight(),
				osmomath.BigDecFromDec(math.LegacyNewDecFromInt(usdcAmount)),
				osmomath.BigDecFromDec(math.LegacyNewDecFromInt(edenAmount)),
			)
			params := k.pk.GetParams(ctx)
			dataLifetime := params.RewardsDataLifetime
			for {
				firstAccum := k.FirstPoolRewardsAccum(ctx, vault.Id)
				if firstAccum.Timestamp == 0 || int64(firstAccum.Timestamp+dataLifetime) >= ctx.BlockTime().Unix() {
					break
				}
				k.DeletePoolRewardsAccum(ctx, firstAccum)
			}
		}
	}
}

// distributeVaultFees handles the logic for distributing management or performance fees.
func (k Keeper) distributeVaultFees(
	ctx sdk.Context,
	vault types.Vault,
	feeRate math.LegacyDec, // can be vault.ManagementFee or shares
	totalBlocksPerYear uint64,
	protocolAddress string,
) {
	var protocolCoins sdk.Coins
	var managerCoins sdk.Coins

	vaultAddress := types.NewVaultAddress(vault.Id)
	coins := k.bk.GetAllBalances(ctx, vaultAddress)
	commitments := k.commitment.GetCommitments(ctx, vaultAddress)

	for _, commitment := range commitments.CommittedTokens {
		if strings.HasPrefix(commitment.Denom, "amm/pool/") {
			poolId, err := GetPoolIdFromShareDenom(commitment.Denom)
			if err != nil {
				k.Logger(ctx).Error("error getting pool id from share denom", "error", err)
				continue
			}
			commitment.Amount = (commitment.Amount.ToLegacyDec().Mul(feeRate).Quo(math.LegacyNewDecFromInt(math.NewInt(int64(totalBlocksPerYear))))).TruncateInt()

			// exit pool
			exitCoins, _, _, _, _, err := k.amm.ExitPool(ctx, vaultAddress, poolId, commitment.Amount, sdk.Coins{}, commitment.Denom, false, true)
			if err != nil {
				k.Logger(ctx).Error("error exiting pool", "error", err)
				continue
			}
			for _, coin := range exitCoins {
				protocolFeeShare := coin.Amount.ToLegacyDec().Mul(vault.ProtocolFeeShare)
				protocolCoins = protocolCoins.Add(sdk.NewCoin(coin.Denom, protocolFeeShare.TruncateInt()))
				coin.Amount = coin.Amount.Sub(protocolFeeShare.TruncateInt())
				managerCoins = managerCoins.Add(sdk.NewCoin(coin.Denom, coin.Amount))
			}
		}
	}

	for _, coin := range coins {
		coin.Amount = (coin.Amount.ToLegacyDec().Mul(feeRate).Quo(math.LegacyNewDecFromInt(math.NewInt(int64(totalBlocksPerYear))))).TruncateInt()

		protocolFeeShare := coin.Amount.ToLegacyDec().Mul(vault.ProtocolFeeShare)
		protocolCoins = protocolCoins.Add(sdk.NewCoin(coin.Denom, protocolFeeShare.TruncateInt()))
		coin.Amount = coin.Amount.Sub(protocolFeeShare.TruncateInt())
		managerCoins = managerCoins.Add(sdk.NewCoin(coin.Denom, coin.Amount))
	}

	// send coins to protocol revenue address and manager address
	err := k.bk.SendCoins(ctx, vaultAddress, sdk.MustAccAddressFromBech32(vault.Manager), managerCoins)
	if err != nil {
		k.Logger(ctx).Error("error sending coins to vault manager", "error", err)
	}
	err = k.bk.SendCoins(ctx, vaultAddress, sdk.MustAccAddressFromBech32(protocolAddress), protocolCoins)
	if err != nil {
		k.Logger(ctx).Error("error sending coins to protocol address", "error", err)
	}
}
