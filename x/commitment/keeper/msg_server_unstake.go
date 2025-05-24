package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/elys-network/elys/v5/x/commitment/types"
	paramtypes "github.com/elys-network/elys/v5/x/parameter/types"
)

func (k msgServer) Unstake(goCtx context.Context, msg *types.MsgUnstake) (*types.MsgUnstakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Asset == paramtypes.Elys {
		if err := k.performUnstakeElys(ctx, msg); err != nil {
			return nil, errorsmod.Wrap(err, "perform elys unstake")
		}
	} else {
		if err := k.performUncommit(ctx, msg); err != nil {
			return nil, errorsmod.Wrap(err, "perform elys uncommit")
		}
	}

	return &types.MsgUnstakeResponse{
		Code:   paramtypes.RES_OK,
		Result: "Unstaking succeed",
	}, nil
}

func (k msgServer) performUnstakeElys(ctx sdk.Context, msg *types.MsgUnstake) error {
	stakingMsgServer := stakingkeeper.NewMsgServerImpl(k.stakingKeeper)

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
	msgMsgUndelegate := stakingtypes.NewMsgUndelegate(address.String(), validator_address.String(), amount)

	if _, err := stakingMsgServer.Undelegate(ctx, msgMsgUndelegate); err != nil { // Discard the response because it's empty
		return errorsmod.Wrap(err, "elys unstake msg")
	}

	return nil
}

func (k msgServer) performUncommit(ctx sdk.Context, msg *types.MsgUnstake) error {
	msgMsgUncommit := types.NewMsgUncommitTokens(msg.Creator, msg.Amount, msg.Asset)

	if err := msgMsgUncommit.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "failed validating msgMsgUncommit")
	}

	_, err := k.UncommitTokens(ctx, msgMsgUncommit) // Discard the response because it's empty
	if err != nil {
		return errorsmod.Wrap(err, "uncommit msg")
	}

	return nil
}
