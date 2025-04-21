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
	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgAddVaultResponse{}, nil
}
