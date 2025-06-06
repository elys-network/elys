package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v6/x/vaults/types"
)

func (k msgServer) UpdateVaultCoins(goCtx context.Context, req *types.MsgUpdateVaultCoins) (*types.MsgUpdateVaultCoinsResponse, error) {
	if k.GetAuthority() != req.Creator {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid creator; expected %s, got %s", k.GetAuthority(), req.Creator)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, types.ErrVaultNotFound
	}

	vault.AllowedCoins = req.AllowedCoins
	vault.RewardCoins = req.RewardCoins

	k.SetVault(ctx, vault)

	return &types.MsgUpdateVaultCoinsResponse{}, nil
}

func (k msgServer) UpdateVaultFees(goCtx context.Context, req *types.MsgUpdateVaultFees) (*types.MsgUpdateVaultFeesResponse, error) {
	if k.GetAuthority() != req.Creator {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid creator; expected %s, got %s", k.GetAuthority(), req.Creator)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, types.ErrVaultNotFound
	}

	vault.ManagementFee = req.ManagementFee
	vault.PerformanceFee = req.PerformanceFee
	vault.ProtocolFeeShare = req.ProtocolFeeShare

	k.SetVault(ctx, vault)

	return &types.MsgUpdateVaultFeesResponse{}, nil
}

func (k msgServer) UpdateVaultLockupPeriod(goCtx context.Context, req *types.MsgUpdateVaultLockupPeriod) (*types.MsgUpdateVaultLockupPeriodResponse, error) {
	if k.GetAuthority() != req.Creator {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid creator; expected %s, got %s", k.GetAuthority(), req.Creator)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, types.ErrVaultNotFound
	}

	vault.LockupPeriod = req.LockupPeriod

	k.SetVault(ctx, vault)

	return &types.MsgUpdateVaultLockupPeriodResponse{}, nil
}

func (k msgServer) UpdateVaultMaxAmountUsd(goCtx context.Context, req *types.MsgUpdateVaultMaxAmountUsd) (*types.MsgUpdateVaultMaxAmountUsdResponse, error) {
	if k.GetAuthority() != req.Creator {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid creator; expected %s, got %s", k.GetAuthority(), req.Creator)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	vault, found := k.GetVault(ctx, req.VaultId)
	if !found {
		return nil, types.ErrVaultNotFound
	}

	vault.MaxAmountUsd = req.MaxAmountUsd

	k.SetVault(ctx, vault)

	return &types.MsgUpdateVaultMaxAmountUsdResponse{}, nil
}
