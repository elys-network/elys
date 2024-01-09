package keeper

import (
	"context"
	"math"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/incentive/types"
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

func (k msgServer) WithdrawRewards(goCtx context.Context, msg *types.MsgWithdrawRewards) (*types.MsgWithdrawRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	// Withdraw rewards
	err = k.ProcessWithdrawRewards(ctx, msg.DelegatorAddress, msg.WithdrawType)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	)
	return &types.MsgWithdrawRewardsResponse{}, nil
}

func (k msgServer) WithdrawValidatorCommission(goCtx context.Context, msg *types.MsgWithdrawValidatorCommission) (*types.MsgWithdrawValidatorCommissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validatorAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	delegator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	found := false
	// Get all delegations
	delegations := k.stk.GetDelegatorDelegations(ctx, delegator, math.MaxUint16)
	for _, del := range delegations {
		// Get validator address
		valAddr := del.GetValidatorAddr()

		// If it is not requested by the validator creator
		if strings.EqualFold(validatorAddr.String(), valAddr.String()) {
			found = true
			break
		}
	}

	// Couldn't find the validator of the delegator.
	if !found {
		return &types.MsgWithdrawValidatorCommissionResponse{}, err
	}

	// Withdraw validator commission
	// Validator will receive commissions from Elys staking only, so the program type is only Rewards_Elys_Program
	// And don't need to input program type
	err = k.RecordWithdrawValidatorCommission(ctx, msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return &types.MsgWithdrawValidatorCommissionResponse{}, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress),
		),
	)

	return &types.MsgWithdrawValidatorCommissionResponse{}, nil
}
