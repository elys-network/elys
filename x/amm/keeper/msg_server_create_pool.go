package keeper

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// CreatePool attempts to create a pool returning the newly created pool ID or an error upon failure.
// The pool creation fee is used to fund the community pool.
// It will create a dedicated module account for the pool and sends the initial liquidity to the created module account.
func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Pay pool creation fee
	params := k.GetParams(ctx)

	if !params.IsCreatorAllowed(msg.Sender) {
		return nil, errors.New("sender is not allowed to create pool")
	}

	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	baseAssetExists := false
	for _, asset := range msg.PoolAssets {
		baseAssetExists = k.CheckBaseAssetExist(ctx, asset.Token.Denom)
		if baseAssetExists {
			break
		}
	}
	if !baseAssetExists {
		return nil, errorsmod.Wrapf(types.ErrOnlyBaseAssetsPoolAllowed, "one of the asset must be from %s", strings.Join(params.BaseAssets, ", "))
	}

	feeAssetExists := k.CheckBaseAssetExist(ctx, msg.PoolParams.FeeDenom)
	if !feeAssetExists {
		return nil, fmt.Errorf("fee denom must be from %s", strings.Join(params.BaseAssets, ", "))
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
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgCreatePoolResponse{
		PoolID: poolId,
	}, nil
}
