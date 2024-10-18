package keeper

import (
	"context"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// CreatePool attempts to create a pool returning the newly created pool ID or an error upon failure.
// The pool creation fee is used to fund the community pool.
// It will create a dedicated module account for the pool and sends the initial liquidity to the created module account.
func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Pay pool creation fee
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	params := k.GetParams(ctx)

	if params.EnableUsdcPairedPoolOnly {
		baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
		if !found {
			return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
		}

		usdc := false
		for _, asset := range msg.PoolAssets {
			if asset.Token.Denom == baseCurrency {
				usdc = true
			}
		}
		if !usdc {
			return nil, errorsmod.Wrapf(types.ErrOnlyUsdcPoolAllowed, "one of the asset must be", ptypes.BaseCurrency)
		}
	}

	if !params.PoolCreationFee.IsNil() && params.PoolCreationFee.IsPositive() {
		feeCoins := sdk.Coins{sdk.NewCoin(ptypes.Elys, params.PoolCreationFee)}
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, feeCoins); err != nil {
			return nil, err
		}
	}

	poolId, err := k.Keeper.CreatePool(ctx, msg)
	if err != nil {
		return &types.MsgCreatePoolResponse{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtPoolCreated,
			sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(poolId, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.GetSigners()[0].String()),
		),
	})

	return &types.MsgCreatePoolResponse{
		PoolID: poolId,
	}, nil
}
