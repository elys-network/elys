package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"

	ammtypes "github.com/elys-network/elys/x/amm/types"
	atypes "github.com/elys-network/elys/x/assetprofile/types"
	tiertypes "github.com/elys-network/elys/x/tier/types"

	"github.com/elys-network/elys/x/vaults/types"
)

func (k msgServer) Deposit(goCtx context.Context, req *types.MsgDeposit) (*types.MsgDepositResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return &types.MsgDepositResponse{}, types.ErrVaultNotFound
	}

	depositer := sdk.MustAccAddressFromBech32(req.Depositor)
	redemptionRate := k.CalculateRedemptionRateForVault(ctx, vault.Id)
	vaultAddress := types.NewVaultAddress(vault.Id)

	depositCoin := sdk.NewCoin(vault.DepositDenom, req.Amount.Amount)
	err := k.bk.SendCoins(ctx, depositer, vaultAddress, sdk.Coins{depositCoin})
	if err != nil {
		return nil, err
	}

	shareDenom := types.GetShareDenomForVault(vault.Id)
	// Initial case
	if redemptionRate.IsZero() {
		redemptionRate = osmomath.OneBigDec()
	}
	shareAmount := (osmomath.BigDecFromSDKInt(depositCoin.Amount).Quo(redemptionRate)).Dec().RoundInt()
	shareCoins := sdk.NewCoins(sdk.NewCoin(shareDenom, shareAmount))

	err = k.bk.MintCoins(ctx, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositer, shareCoins)
	if err != nil {
		return nil, err
	}

	_, found = k.assetProfileKeeper.GetEntry(ctx, shareDenom)
	if !found {
		// Set an entity to assetprofile
		entry := atypes.Entry{
			Authority:                vaultAddress.String(),
			BaseDenom:                shareDenom,
			Decimals:                 6, // TODO: Get from assetprofile of deposit denom or keep it as 6
			Denom:                    shareDenom,
			Path:                     "",
			IbcChannelId:             "",
			IbcCounterpartyChannelId: "",
			DisplayName:              shareDenom,
			DisplaySymbol:            "",
			Network:                  "",
			Address:                  "",
			ExternalSymbol:           "",
			TransferLimit:            "",
			Permissions:              make([]string, 0),
			UnitDenom:                "",
			IbcCounterpartyDenom:     "",
			IbcCounterpartyChainId:   "",
			CommitEnabled:            true,
			WithdrawEnabled:          true,
		}

		k.assetProfileKeeper.SetEntry(ctx, entry)
	}

	// Commit LP token
	lockUntil := uint64(ctx.BlockTime().Second()) + vault.LockupPeriod
	err = k.commitement.CommitLiquidTokens(ctx, depositer, shareDenom, shareAmount, lockUntil)
	if err != nil {
		return nil, err
	}

	// convert input amount to USD value
	usdValue := k.amm.CalculateUSDValue(ctx, vault.DepositDenom, req.Amount.Amount)
	userData, found := k.GetUserData(ctx, depositer.String())
	if found {
		userData.PrincipalAmount = userData.PrincipalAmount.Add(usdValue.Dec())
	} else {
		userData = types.UserData{
			User:            depositer.String(),
			PrincipalAmount: usdValue.Dec(),
		}
	}

	k.SetUserData(ctx, userData)

	return &types.MsgDepositResponse{
		VaultId: vault.Id,
		Shares:  shareAmount,
	}, nil
}

func (k Keeper) VaultUsdValue(ctx sdk.Context, vaultId uint64) (osmomath.BigDec, error) {
	vault, found := k.GetVault(ctx, vaultId)
	if !found {
		return osmomath.ZeroBigDec(), types.ErrVaultNotFound
	}
	vaultAddress := types.NewVaultAddress(vaultId)
	totalValue := osmomath.ZeroBigDec()
	commitments := k.commitement.GetCommitments(ctx, vaultAddress)
	// TODO: Handle zero values for denom, we should issue shares if price is not available
	for _, commitment := range commitments.CommittedTokens {
		// Pool balance
		if strings.HasPrefix(commitment.Denom, "amm/pool") {
			poolId, err := ammtypes.GetPoolIdFromShareDenom(commitment.Denom)
			if err != nil {
				continue
			}
			pool, found := k.amm.GetPool(ctx, poolId)
			if !found {
				continue
			}
			info := k.amm.PoolExtraInfo(ctx, pool, tiertypes.OneDay)
			amount := osmomath.BigDecFromSDKInt(commitment.Amount)
			totalValue = totalValue.Add(amount.Mul(osmomath.BigDecFromDec(info.LpTokenPrice)).Quo(osmomath.BigDecFromSDKInt(ammtypes.OneShare)))
		}
	}
	for _, coin := range vault.AllowedCoins {
		if !strings.HasPrefix(coin, "amm/pool") {
			balance := k.bk.GetBalance(ctx, vaultAddress, coin)
			totalValue = totalValue.Add(k.amm.CalculateUSDValue(ctx, coin, balance.Amount))
		}
	}
	for _, coin := range vault.RewardCoins {
		if !strings.HasPrefix(coin, "amm/pool") {
			balance := k.bk.GetBalance(ctx, vaultAddress, coin)
			totalValue = totalValue.Add(k.amm.CalculateUSDValue(ctx, coin, balance.Amount))
		}
	}
	return totalValue, nil
}

func (k Keeper) CalculateRedemptionRateForVault(ctx sdk.Context, vaultId uint64) osmomath.BigDec {
	totalShares := k.bk.GetSupply(ctx, types.GetShareDenomForVault(vaultId))

	if totalShares.Amount.IsZero() {
		return osmomath.ZeroBigDec()
	}

	// TODO: Handle zero values for denom, we should not issue shares if price is not available
	// TODO: Should it be based on deposit denom value ?
	usdValue, err := k.VaultUsdValue(ctx, vaultId)
	if err != nil {
		return osmomath.ZeroBigDec()
	}

	return usdValue.Quo(osmomath.BigDecFromSDKInt(totalShares.Amount))
}
