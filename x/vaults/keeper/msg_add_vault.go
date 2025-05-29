package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v5/x/vaults/types"
)

func (k msgServer) AddVault(goCtx context.Context, req *types.MsgAddVault) (*types.MsgAddVaultResponse, error) {
	if k.GetAuthority() != req.Creator {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Creator)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	vaultId := k.GetNextVaultId(ctx)
	vault := types.Vault{
		Id:                    vaultId,
		DepositDenom:          req.DepositDenom,
		MaxAmountUsd:          req.MaxAmountUsd,
		AllowedCoins:          req.AllowedCoins,
		RewardCoins:           req.RewardCoins,
		Manager:               req.Manager,
		ManagementFee:         req.ManagementFee,
		PerformanceFee:        req.PerformanceFee,
		BenchmarkCoin:         req.BenchmarkCoin,
		ProtocolFeeShare:      req.ProtocolFeeShare,
		LockupPeriod:          req.LockupPeriod,
		WithdrawalUsdValue:    math.LegacyZeroDec(),
		SumOfDepositsUsdValue: math.LegacyZeroDec(),
	}

	k.SetVault(ctx, vault)

	vaultAddress := types.NewVaultAddress(vault.Id)
	vaultName := types.GetVaultIdModuleName(vault.Id)
	if err := types.CreateModuleAccount(ctx, k.accountKeeper, vaultAddress, vaultName); err != nil {
		return &types.MsgAddVaultResponse{}, fmt.Errorf("creating vault module account for id %d: %w", vault.Id, err)
	}

	return &types.MsgAddVaultResponse{}, nil
}
