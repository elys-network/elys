package keeper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	tiertypes "github.com/elys-network/elys/v7/x/tier/types"
	"github.com/elys-network/elys/v7/x/vaults/types"
)

func (k Keeper) Vault(goCtx context.Context, req *types.QueryVaultRequest) (*types.QueryVaultResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	vaultAndData, err := k.GetVaultAndData(ctx, req.VaultId, req.Days)
	if err != nil {
		return nil, err
	}

	return &types.QueryVaultResponse{Vault: vaultAndData}, nil
}

func (k Keeper) Vaults(goCtx context.Context, req *types.QueryVaultsRequest) (*types.QueryVaultsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	vaults := k.GetAllVaults(ctx)
	vaultsAndData := []types.VaultAndData{}
	for _, vault := range vaults {
		vaultAndData, err := k.GetVaultAndData(ctx, vault.Id, req.Days)
		if err != nil {
			return nil, err
		}
		vaultsAndData = append(vaultsAndData, vaultAndData)
	}

	return &types.QueryVaultsResponse{Vaults: vaultsAndData}, nil
}

func (k Keeper) GetVaultAndData(ctx sdk.Context, vaultId uint64, days uint64) (types.VaultAndData, error) {
	vault, found := k.GetVault(ctx, vaultId)
	if !found {
		return types.VaultAndData{}, fmt.Errorf("vault not found")
	}

	edenApr := k.EdenApr(ctx, vaultId)
	pnlUsd := k.GetPnlTotal(ctx, strconv.FormatUint(vaultId, 10), int(days))
	totalDepositsUsd, _ := k.VaultUsdValue(ctx, vaultId)
	// Deposit denom usd value
	balance := k.bk.GetBalance(ctx, types.NewVaultAddress(vaultId), vault.DepositDenom)
	depositDenomUsdValue := k.amm.CalculateUSDValue(ctx, vault.DepositDenom, balance.Amount)
	var depositsUsed osmomath.BigDec
	if totalDepositsUsd.Equal(depositDenomUsdValue) {
		depositsUsed = osmomath.OneBigDec()
	} else {
		depositsUsed = depositDenomUsdValue.Quo(totalDepositsUsd.Sub(depositDenomUsdValue))
	}
	positions, err := k.GetVaultPositions(ctx, vaultId)
	if err != nil {
		return types.VaultAndData{}, err
	}

	return types.VaultAndData{
		Vault:            &vault,
		EdenApr:          edenApr.Dec(),
		PnlUsd:           pnlUsd.Dec(),
		TotalDepositsUsd: totalDepositsUsd.Dec(),
		DepositsUsed:     depositsUsed.Dec(),
		Positions:        positions,
	}, nil
}

func (k Keeper) VaultValue(goCtx context.Context, req *types.QueryVaultValue) (*types.QueryVaultValueResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, fmt.Errorf("vault not found")
	}

	usdValue, err := k.VaultUsdValue(ctx, req.VaultId)
	if err != nil {
		return nil, err
	}

	return &types.QueryVaultValueResponse{UsdValue: usdValue.Dec()}, nil
}

func (k Keeper) VaultPositions(goCtx context.Context, req *types.QueryVaultPositionsRequest) (*types.QueryVaultPositionsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	positions, err := k.GetVaultPositions(ctx, req.VaultId)
	if err != nil {
		return nil, err
	}
	return &types.QueryVaultPositionsResponse{Positions: positions}, nil
}

func (k Keeper) GetVaultPositions(ctx sdk.Context, vaultId uint64) ([]types.PositionToken, error) {
	vaultAddress := types.NewVaultAddress(vaultId)
	positions := []types.PositionToken{}
	commitments := k.commitment.GetCommitments(ctx, vaultAddress)
	for _, commitment := range commitments.CommittedTokens {
		if strings.HasPrefix(commitment.Denom, "amm/pool") {
			poolId, err := ammtypes.GetPoolIdFromShareDenom(commitment.Denom)
			if err != nil {
				return nil, fmt.Errorf("invalid pool denom: %s", commitment.Denom)
			}
			pool, found := k.amm.GetPool(ctx, poolId)
			if !found {
				return nil, fmt.Errorf("pool not found for denom: %s", commitment.Denom)
			}
			info := k.amm.PoolExtraInfo(ctx, pool, tiertypes.OneDay)
			amount := osmomath.BigDecFromSDKInt(commitment.Amount)
			if info.LpTokenPrice.IsZero() {
				return nil, fmt.Errorf("no price available for pool denom: %s", commitment.Denom)
			}
			token := types.PositionToken{
				TokenDenom:    commitment.Denom,
				TokenAmount:   amount.Dec(),
				TokenUsdValue: amount.Mul(osmomath.BigDecFromDec(info.LpTokenPrice)).Quo(osmomath.BigDecFromSDKInt(ammtypes.OneShare)).Dec(),
			}
			positions = append(positions, token)
		}
	}
	balances := k.bk.GetAllBalances(ctx, vaultAddress)
	for _, balance := range balances {
		usdVal := k.amm.CalculateUSDValue(ctx, balance.Denom, balance.Amount)
		if usdVal.IsZero() {
			return nil, fmt.Errorf("no price available for denom: %s", balance.Denom)
		}
		token := types.PositionToken{
			TokenDenom:    balance.Denom,
			TokenAmount:   osmomath.BigDecFromSDKInt(balance.Amount).Dec(),
			TokenUsdValue: usdVal.Dec(),
		}
		positions = append(positions, token)
	}
	return positions, nil
}

func (k Keeper) EdenApr(ctx sdk.Context, vaultId uint64) osmomath.BigDec {
	var edenApr osmomath.BigDec
	totalBlocksPerYear := k.pk.GetParams(ctx).TotalBlocksPerYear
	usdcDenomPrice := k.oracleKeeper.GetDenomPrice(ctx, ptypes.BaseCurrency)

	tvl, err := k.VaultUsdValue(ctx, vaultId)
	if err != nil {
		return osmomath.ZeroBigDec()
	}

	firstAccum := k.FirstPoolRewardsAccum(ctx, vaultId)
	lastAccum := k.LastPoolRewardsAccum(ctx, vaultId)
	if lastAccum.Timestamp == 0 || tvl.IsZero() {
		return osmomath.ZeroBigDec()
	}

	if firstAccum.Timestamp == lastAccum.Timestamp {
		edenApr = osmomath.BigDecFromDec(lastAccum.EdenReward.
			Mul(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(int64(totalBlocksPerYear)))).
			Mul(usdcDenomPrice.Dec()).
			Quo(tvl.Dec()))
	} else {
		duration := lastAccum.Timestamp - firstAccum.Timestamp
		secondsInYear := int64(86400 * 360)

		edenApr = osmomath.BigDecFromDec(lastAccum.EdenReward.Sub(firstAccum.EdenReward).
			Mul(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(secondsInYear))).
			Quo(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(int64(duration)))).
			Mul(usdcDenomPrice.Dec()).
			Quo(tvl.Dec()))
	}
	return edenApr
}

// func (k Keeper) PnL(goCtx context.Context, req *types.QueryPnLRequest) (*types.QueryPnLResponse, error) {
// 	if req == nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid request")
// 	}
// 	ctx := sdk.UnwrapSDKContext(goCtx)

// 	vaults := k.GetAllVaults(ctx)
// 	pnls := []types.PnlResponse{}
// 	for _, vault := range vaults {
// 		userData, _ := k.GetUserData(ctx, req.Address, vault.Id)
// 		// get vault usd value
// 		commitments := k.commitment.GetCommitments(ctx, types.NewVaultAddress(vault.Id))
// 		shareAmount := sdkmath.ZeroInt()
// 		for _, commitment := range commitments.CommittedTokens {
// 			if commitment.Denom == types.GetShareDenomForVault(vault.Id) {
// 				shareAmount = shareAmount.Add(commitment.Amount)
// 			}
// 		}
// 		vaultValue, err := k.VaultUsdValue(ctx, vault.Id)
// 		if err != nil {
// 			return nil, err
// 		}
// 		totalShares := k.bk.GetSupply(ctx, types.GetShareDenomForVault(vault.Id)).Amount

// 		shareUsdValue := vaultValue.Dec().Mul(shareAmount.ToLegacyDec().Quo(sdkmath.LegacyNewDecFromInt(totalShares)))
// 		currentBalanceUsd := userData.TotalDepositsUsd.Sub(userData.TotalWithdrawalsUsd)
// 		profitAndLossUsd := shareUsdValue.Sub(currentBalanceUsd)

// 		unclaimedEden, found := k.GetUserRewardInfo(ctx, sdk.MustAccAddressFromBech32(req.Address), vault.Id, ptypes.Eden)
// 		if found && !unclaimedEden.RewardPending.IsZero() {
// 			edenUsdValue := k.amm.CalculateUSDValue(ctx, ptypes.Eden, unclaimedEden.RewardPending.TruncateInt())
// 			userData.EdenUsdValue = userData.EdenUsdValue.Add(edenUsdValue.Dec())
// 			userData.EdenAmount = userData.EdenAmount.Add(unclaimedEden.RewardPending.TruncateInt())
// 		}

// 		pnls = append(pnls, types.PnlResponse{
// 			PnlUsd:            profitAndLossUsd,
// 			EdenUsdValue:      userData.EdenUsdValue,
// 			CurrentBalanceUsd: currentBalanceUsd,
// 			VaultId:           vault.Id,
// 			EdenAmount:        userData.EdenAmount,
// 		})
// 	}

// 	return &types.QueryPnLResponse{
// 		Pnls: pnls,
// 	}, nil
// }

// func (k Keeper) PnlApr(ctx sdk.Context, vaultId uint64) osmomath.BigDec {
// 	vault, found := k.GetVault(ctx, vaultId)
// 	if !found {
// 		return osmomath.ZeroBigDec()
// 	}
// 	pnlApr := k.GetPnlApr(ctx, vaultId)
// 	return pnlApr
// }
