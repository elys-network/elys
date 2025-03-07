package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/elys-network/elys/x/commitment/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (k msgServer) Stake(goCtx context.Context, msg *types.MsgStake) (*types.MsgStakeResponse, error) {
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
		return errors.New("invalid amount")
	}

	// Don't allow vested tokens to be staked
	// Retrieve the delegator account
	delegatorAcc := k.accountKeeper.GetAccount(ctx, address)
	if _, ok := delegatorAcc.(banktypes.VestingAccount); ok {
		spendableCoins := k.bankKeeper.SpendableCoins(ctx, address)
		if msg.Amount.GT(spendableCoins.AmountOf(msg.Asset)) {
			return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "cannot delegate vested tokens")
		}
	}

	msgMsgDelegate := stakingtypes.NewMsgDelegate(address.String(), validator_address.String(), amount)

	if _, err := stakingMsgServer.Delegate(ctx, msgMsgDelegate); err != nil { // Discard the response because it's empty
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
