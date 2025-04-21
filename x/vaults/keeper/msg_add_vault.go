package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/vaults/types"
)

func (k msgServer) AddVault(goCtx context.Context, req *types.MsgAddVault) (*types.MsgAddVaultResponse, error) {
	if k.GetAuthority() != req.Creator {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Creator)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	vaultId := k.GetNextVaultId(ctx)
	vault := types.Vault{
		Id:             vaultId,
		DepositDenom:   req.DepositDenom,
		MaxAmountUsd:   req.MaxAmountUsd,
		AllowedCoins:   req.AllowedCoins,
		AllowedActions: req.AllowedActions,
	}
	k.SetVault(ctx, vault)

	return &types.MsgAddVaultResponse{}, nil
}
