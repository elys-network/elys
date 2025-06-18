package keeper

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"

	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	atypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	tiertypes "github.com/elys-network/elys/v6/x/tier/types"

	"github.com/elys-network/elys/v6/x/vaults/types"
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

	if req.Amount.Denom != vault.DepositDenom {
		return nil, types.ErrInvalidDepositDenom
	}

	k.DeductPerformanceFee(ctx)

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

	// Ensure asset profile entry exists for share denom
	err = k.Keeper.EnsureAssetProfileEntry(ctx, shareDenom, vaultAddress.String())
	if err != nil {
		return nil, err
	}

	// Commit LP token
	lockUntil := uint64(ctx.BlockTime().Second()) + vault.LockupPeriod
	err = k.commitment.CommitLiquidTokens(ctx, depositer, shareDenom, shareAmount, lockUntil)
	if err != nil {
		return nil, err
	}

	// Set sum of deposits usd value
	usdValue := k.amm.CalculateUSDValue(ctx, req.Amount.Denom, req.Amount.Amount)
	if usdValue.IsZero() {
		return nil, types.ErrDepositValueZero
	}
	vault.SumOfDepositsUsdValue = vault.SumOfDepositsUsdValue.Add(usdValue.Dec())
	k.SetVault(ctx, vault)

	k.AfterDeposit(ctx, vault.Id, depositer, shareAmount)

	return &types.MsgDepositResponse{
		VaultId: vault.Id,
		Shares:  shareAmount,
	}, nil
}

func (k Keeper) VaultUsdValue(ctx sdk.Context, vaultId uint64) (osmomath.BigDec, error) {
	vaultAddress := types.NewVaultAddress(vaultId)
	totalValue := osmomath.ZeroBigDec()
	commitments := k.commitment.GetCommitments(ctx, vaultAddress)
	for _, commitment := range commitments.CommittedTokens {
		if strings.HasPrefix(commitment.Denom, "amm/pool") {
			poolId, err := ammtypes.GetPoolIdFromShareDenom(commitment.Denom)
			if err != nil {
				return osmomath.ZeroBigDec(), fmt.Errorf("invalid pool denom: %s", commitment.Denom)
			}
			pool, found := k.amm.GetPool(ctx, poolId)
			if !found {
				return osmomath.ZeroBigDec(), fmt.Errorf("pool not found for denom: %s", commitment.Denom)
			}
			info := k.amm.PoolExtraInfo(ctx, pool, tiertypes.OneDay)
			amount := osmomath.BigDecFromSDKInt(commitment.Amount)
			if info.LpTokenPrice.IsZero() {
				return osmomath.ZeroBigDec(), fmt.Errorf("no price available for pool denom: %s", commitment.Denom)
			}
			totalValue = totalValue.Add(amount.Mul(osmomath.BigDecFromDec(info.LpTokenPrice)).Quo(osmomath.BigDecFromSDKInt(ammtypes.OneShare)))
		}
	}
	balances := k.bk.GetAllBalances(ctx, vaultAddress)
	for _, balance := range balances {
		usdVal := k.amm.CalculateUSDValue(ctx, balance.Denom, balance.Amount)
		if usdVal.IsZero() {
			return osmomath.ZeroBigDec(), fmt.Errorf("no price available for denom: %s", balance.Denom)
		}
		totalValue = totalValue.Add(usdVal)
	}
	return totalValue, nil
}

func (k Keeper) CalculateRedemptionRateForVault(ctx sdk.Context, vaultId uint64) osmomath.BigDec {
	totalShares := k.bk.GetSupply(ctx, types.GetShareDenomForVault(vaultId))

	if totalShares.Amount.IsZero() {
		return osmomath.ZeroBigDec()
	}

	usdValue, err := k.VaultUsdValue(ctx, vaultId)
	if err != nil || usdValue.IsZero() {
		return osmomath.ZeroBigDec()
	}
	return usdValue.Quo(osmomath.BigDecFromSDKInt(totalShares.Amount))
}

// EnsureAssetProfileEntry creates an asset profile entry if it does not exist
func (k Keeper) EnsureAssetProfileEntry(ctx sdk.Context, denom string, authority string) error {
	_, found := k.assetProfileKeeper.GetEntry(ctx, denom)
	if found {
		return nil
	}
	entry := atypes.Entry{
		Authority:                authority,
		BaseDenom:                denom,
		Decimals:                 6, // TODO: Get from assetprofile of deposit denom or keep it as 6
		Denom:                    denom,
		Path:                     "",
		IbcChannelId:             "",
		IbcCounterpartyChannelId: "",
		DisplayName:              denom,
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
	return nil
}
