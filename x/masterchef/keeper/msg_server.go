package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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

func (k msgServer) AddExternalRewardDenom(goCtx context.Context, msg *types.MsgAddExternalRewardDenom) (*types.MsgAddExternalRewardDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	params := k.GetParams(ctx)

	index := -1
	for i, rewardDenom := range params.SupportedRewardDenoms {
		if rewardDenom.Denom == msg.RewardDenom {
			index = i
			break
		}
	}

	if index == -1 && msg.Supported {
		params.SupportedRewardDenoms = append(params.SupportedRewardDenoms, &types.SupportedRewardDenom{
			Denom:     msg.RewardDenom,
			MinAmount: msg.MinAmount,
		})
	}

	if index != -1 && !msg.Supported {
		params.SupportedRewardDenoms = append(
			params.SupportedRewardDenoms[:index],
			params.SupportedRewardDenoms[index+1:]...,
		)
	}

	if index != -1 && msg.Supported {
		params.SupportedRewardDenoms[index].MinAmount = msg.MinAmount
	}

	k.SetParams(ctx, params)

	return &types.MsgAddExternalRewardDenomResponse{}, nil
}

func (k msgServer) AddExternalIncentive(goCtx context.Context, msg *types.MsgAddExternalIncentive) (*types.MsgAddExternalIncentiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	if msg.FromBlock < uint64(ctx.BlockHeight()) {
		return nil, status.Error(codes.InvalidArgument, "invalid from block")
	}
	if msg.FromBlock >= msg.ToBlock {
		return nil, status.Error(codes.InvalidArgument, "invalid block range")
	}
	if msg.AmountPerBlock.IsZero() {
		return nil, status.Error(codes.InvalidArgument, "invalid amount per block")
	}

	found := false
	params := k.GetParams(ctx)
	for _, rewardDenom := range params.SupportedRewardDenoms {
		if msg.RewardDenom == rewardDenom.Denom {
			found = true
			if msg.AmountPerBlock.Mul(
				sdk.NewInt(int64(msg.ToBlock - msg.FromBlock)),
			).LT(rewardDenom.MinAmount) {
				return nil, status.Error(codes.InvalidArgument, "too small amount")
			}
			break
		}
	}
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid reward denom")
	}

	amount := msg.AmountPerBlock.Mul(sdk.NewInt(int64(msg.ToBlock - msg.FromBlock)))
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.Coins{sdk.NewCoin(msg.RewardDenom, amount)})
	if err != nil {
		return nil, err
	}

	k.Keeper.AddExternalIncentive(ctx, types.ExternalIncentive{
		Id:             0,
		RewardDenom:    msg.RewardDenom,
		PoolId:         msg.PoolId,
		FromBlock:      msg.FromBlock,
		ToBlock:        msg.ToBlock,
		AmountPerBlock: msg.AmountPerBlock,
		Apr:            math.LegacyZeroDec(),
	})

	return &types.MsgAddExternalIncentiveResponse{}, nil
}

func (k Keeper) ClaimRewards(ctx sdk.Context, sender sdk.AccAddress, poolIds []uint64) error {
	coins := sdk.NewCoins()
	for _, poolId := range poolIds {
		k.AfterWithdraw(ctx, poolId, sender.String(), sdk.ZeroInt())

		for _, rewardDenom := range k.GetRewardDenoms(ctx, poolId) {
			userRewardInfo, found := k.GetUserRewardInfo(ctx, sender.String(), poolId, rewardDenom)
			if found && userRewardInfo.RewardPending.IsPositive() {
				coin := sdk.NewCoin(rewardDenom, userRewardInfo.RewardPending.TruncateInt())
				coins = coins.Add(coin)

				userRewardInfo.RewardPending = sdk.ZeroDec()
				k.SetUserRewardInfo(ctx, userRewardInfo)
			}
		}
	}

	// Transfer rewards (Eden/EdenB is transferred through commitment module)
	err := k.cmk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, coins)
	if err != nil {
		return err
	}
	return nil
}

func (k msgServer) ClaimRewards(goCtx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	if len(msg.PoolIds) == 0 {
		allPools := k.GetAllPools(ctx)
		for _, pool := range allPools {
			msg.PoolIds = append(msg.PoolIds, pool.PoolId)
		}
	}

	err := k.Keeper.ClaimRewards(ctx, sender, msg.PoolIds)
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimRewardsResponse{}, nil
}

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	k.SetParams(ctx, msg.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}

func (k msgServer) UpdatePoolMultipliers(goCtx context.Context, msg *types.MsgUpdatePoolMultipliers) (*types.MsgUpdatePoolMultipliersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	k.Keeper.UpdatePoolMultipliers(ctx, msg.PoolMultipliers)

	return &types.MsgUpdatePoolMultipliersResponse{}, nil
}
