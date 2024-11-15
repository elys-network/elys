package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/elys-network/elys/x/commitment/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (k msgServer) Stake(goCtx context.Context, msg *types.MsgStake) (*types.MsgStakeResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Asset == paramtypes.Elys {
		if err := k.performStakeElys(ctx, msg); err != nil {
			return nil, errorsmod.Wrap(err, "perform elys stake")
		}
	} else {
		if err := k.performCommit(ctx, msg); err != nil {
			return nil, errorsmod.Wrap(err, "perform elys commit")
		}
	}

	return &types.MsgStakeResponse{
		Code:   paramtypes.RES_OK,
		Result: "Staking succeed",
	}, nil
}

func (k msgServer) performStakeElys(ctx sdk.Context, msg *types.MsgStake) error {
	stakingKeeper, ok := k.stakingKeeper.(*stakingkeeper.Keeper)
	if !ok {
		return errorsmod.Wrap(errorsmod.Error{}, "staking keeper")
	}

	msgServer := stakingkeeper.NewMsgServerImpl(stakingKeeper)

	address, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrap(err, "invalid address")
	}

	validator_address, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return errorsmod.Wrap(err, "invalid validator address")
	}

	amount := sdk.NewCoin(msg.Asset, msg.Amount)
	if !amount.IsValid() || amount.Amount.IsZero() {
		return fmt.Errorf("invalid amount")
	}
	msgMsgDelegate := stakingtypes.NewMsgDelegate(address.String(), validator_address.String(), amount)

	if _, err := msgServer.Delegate(ctx, msgMsgDelegate); err != nil { // Discard the response because it's empty
		return errorsmod.Wrap(err, "elys stake msg")
	}

	return nil
}

func (k msgServer) performCommit(ctx sdk.Context, msg *types.MsgStake) error {
	msgMsgCommit := types.NewMsgCommitClaimedRewards(msg.Creator, msg.Amount, msg.Asset)

	if err := msgMsgCommit.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "failed validating msgMsgCommit")
	}

	_, err := k.CommitClaimedRewards(ctx, msgMsgCommit) // Discard the response because it's empty
	if err != nil {
		return errorsmod.Wrap(err, "commit msg")
	}

	return nil
}
