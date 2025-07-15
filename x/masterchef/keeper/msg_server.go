package keeper

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v6/x/masterchef/types"
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
		params.SupportedRewardDenoms = slices.Delete(
			params.SupportedRewardDenoms, index,
			index+1,
		)
	}

	if index != -1 && msg.Supported {
		params.SupportedRewardDenoms[index].MinAmount = msg.MinAmount
	}

	k.SetParams(ctx, params)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtAddExternalRewardDenom,
			sdk.NewAttribute(types.AttributeRewardDenom, msg.RewardDenom),
			sdk.NewAttribute(types.AttributeMinAmount, msg.MinAmount.String()),
			sdk.NewAttribute(types.AttributeSupported, fmt.Sprintf("%t", msg.Supported)),
		),
	})
	return &types.MsgAddExternalRewardDenomResponse{}, nil
}

func (k msgServer) AddExternalIncentive(goCtx context.Context, msg *types.MsgAddExternalIncentive) (*types.MsgAddExternalIncentiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	if msg.FromBlock < ctx.BlockHeight() {
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
			if msg.AmountPerBlock.Mul(math.NewInt(msg.ToBlock - msg.FromBlock)).LT(rewardDenom.MinAmount) {
				return nil, status.Error(codes.InvalidArgument, "too small amount")
			}
			break
		}
	}
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid reward denom")
	}

	amount := msg.AmountPerBlock.Mul(math.NewInt(msg.ToBlock - msg.FromBlock))
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.Coins{sdk.NewCoin(msg.RewardDenom, amount)})
	if err != nil {
		return nil, err
	}

	externalIncentive := types.ExternalIncentive{
		Id:             k.GetExternalIncentiveIndex(ctx),
		RewardDenom:    msg.RewardDenom,
		PoolId:         msg.PoolId,
		FromBlock:      msg.FromBlock,
		ToBlock:        msg.ToBlock,
		AmountPerBlock: msg.AmountPerBlock,
		Apr:            math.LegacyZeroDec(),
	}
	k.Keeper.SetExternalIncentive(ctx, externalIncentive)
	k.SetExternalIncentiveIndex(ctx, externalIncentive.Id+1)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtAddExternalIncentive,
			sdk.NewAttribute(types.AttributeRewardDenom, msg.RewardDenom),
			sdk.NewAttribute(types.AttributePoolId, fmt.Sprintf("%d", msg.PoolId)),
			sdk.NewAttribute(types.AttributeFromBlock, fmt.Sprintf("%d", msg.FromBlock)),
			sdk.NewAttribute(types.AttributeToBlock, fmt.Sprintf("%d", msg.ToBlock)),
			sdk.NewAttribute(types.AttributeAmountPerBlock, fmt.Sprintf("%d", msg.AmountPerBlock)),
		),
	})
	return &types.MsgAddExternalIncentiveResponse{}, nil
}

func (k Keeper) ClaimRewards(ctx sdk.Context, sender sdk.AccAddress, poolIds []uint64, recipient sdk.AccAddress) error {
	coins := sdk.NewCoins()
	rewardPoolIds := []string{}
	for _, poolId := range poolIds {
		k.AfterWithdraw(ctx, poolId, sender, math.ZeroInt())

		for _, rewardDenom := range k.GetRewardDenoms(ctx, poolId) {
			userRewardInfo, found := k.GetUserRewardInfo(ctx, sender, poolId, rewardDenom)
			if found && userRewardInfo.RewardPending.IsPositive() {
				coin := sdk.NewCoin(rewardDenom, userRewardInfo.RewardPending.TruncateInt())
				coins = coins.Add(coin)
				rewardPoolIds = append(rewardPoolIds, strconv.FormatUint(poolId, 10))

				userRewardInfo.RewardPending = math.LegacyZeroDec()
				if userRewardInfo.RewardDebt.IsZero() {
					k.RemoveUserRewardInfo(ctx, userRewardInfo.GetUserAccount(), userRewardInfo.PoolId, userRewardInfo.RewardDenom)
				} else {
					k.SetUserRewardInfo(ctx, userRewardInfo)
				}
			}
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtClaimRewards,
			sdk.NewAttribute(types.AttributeSender, sender.String()),
			sdk.NewAttribute(types.AttributeRecipient, recipient.String()),
			sdk.NewAttribute(types.AttributePoolIds, strings.Join(rewardPoolIds, ",")),
			sdk.NewAttribute(sdk.AttributeKeyAmount, coins.String()),
		),
	})

	// Transfer rewards (Eden/EdenB is transferred through commitment module)
	err := k.commitmentKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins)
	if err != nil {
		return err
	}

	return nil
}

func (k msgServer) ClaimRewards(goCtx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	if len(msg.PoolIds) == 0 {
		allPools := k.GetAllPoolInfos(ctx)
		for _, pool := range allPools {
			msg.PoolIds = append(msg.PoolIds, pool.PoolId)
		}
	}

	err := k.Keeper.ClaimRewards(ctx, sender, msg.PoolIds, sender)
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

	if k.CheckBlockedAddress(msg.Params) {
		return nil, fmt.Errorf("protocol revenue address is blocked")
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

func (k msgServer) ToggleTakerFeeSwapAndBurn(goCtx context.Context, msg *types.MsgToggleTakerFeeSwapAndBurn) (*types.MsgToggleTakerFeeSwapAndBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	takerManager := k.GetParams(ctx).TakerManager

	if takerManager != msg.Sender {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", takerManager, msg.Sender)
	}

	k.ProcessTakerFee(ctx)

	return &types.MsgToggleTakerFeeSwapAndBurnResponse{}, nil
}
