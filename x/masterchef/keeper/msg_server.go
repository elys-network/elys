package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) AddExternalIncentive(goCtx context.Context, msg *types.MsgAddExternalIncentive) (*types.MsgAddExternalIncentiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.FromBlock < uint64(ctx.BlockHeight()) {
		return nil, status.Error(codes.InvalidArgument, "invalid from block")
	}
	if msg.FromBlock >= msg.ToBlock {
		return nil, status.Error(codes.InvalidArgument, "invalid block range")
	}
	if msg.AmountPerBlock.IsZero() {
		return nil, status.Error(codes.InvalidArgument, "invalid amount per block")
	}

	err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		sdk.MustAccAddressFromBech32(msg.Sender),
		types.ModuleName,
		sdk.NewCoins(
			sdk.NewCoin(
				msg.RewardDenom,
				msg.AmountPerBlock.Mul(
					sdk.NewInt(int64(msg.ToBlock-msg.FromBlock)),
				),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	k.SetExternalIncentive(
		ctx,
		types.ExternalIncentive{
			Id:             0,
			RewardDenom:    msg.RewardDenom,
			PoolId:         msg.PoolId,
			FromBlock:      msg.FromBlock,
			ToBlock:        msg.ToBlock,
			AmountPerBlock: msg.AmountPerBlock,
		},
	)

	return &types.MsgAddExternalIncentiveResponse{}, nil
}

func (k msgServer) ClaimRewards(goCtx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	coins := sdk.NewCoins()

	for _, poolId := range msg.PoolIds {
		k.AfterWithdraw(
			ctx,
			poolId,
			msg.Sender,
			sdk.ZeroInt(),
		)

		for _, rewardDenom := range k.GetRewardDenoms(ctx, poolId) {
			userRewardInfo, found := k.GetUserRewardInfo(ctx, msg.Sender, poolId, rewardDenom)

			if found && userRewardInfo.RewardPending.IsPositive() {
				coins = coins.Add(
					sdk.NewCoin(
						rewardDenom,
						userRewardInfo.RewardPending.TruncateInt(),
					),
				)

				userRewardInfo.RewardPending = sdk.ZeroDec()
				k.SetUserRewardInfo(ctx, userRewardInfo)
			}
		}
	}

	// TODO: consider Eden
	err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		sdk.MustAccAddressFromBech32(msg.Sender),
		coins,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimRewardsResponse{}, nil
}
