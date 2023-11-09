package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/elys-network/elys/x/commitment/client/wasm/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (m *Messenger) msgStake(ctx sdk.Context, contractAddr sdk.AccAddress, msgStake *types.MsgStake) ([]sdk.Event, [][]byte, error) {
	var res *types.RequestResponse
	var err error
	if msgStake.Asset == paramtypes.Elys {
		res, err = performMsgStakeElys(m.stakingKeeper, ctx, contractAddr, msgStake)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "perform elys stake")
		}
	} else {
		res, err = performMsgCommit(m.keeper, ctx, contractAddr, msgStake)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "perform elys stake")
		}
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize stake")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func performMsgStakeElys(f *stakingkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgStake *types.MsgStake) (*types.RequestResponse, error) {
	if msgStake == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid staking parameter"}
	}

	msgServer := stakingkeeper.NewMsgServerImpl(f)
	address, err := sdk.AccAddressFromBech32(msgStake.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	validator_address, err := sdk.ValAddressFromBech32(msgStake.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	amount := sdk.NewCoin(msgStake.Asset, msgStake.Amount)
	msgMsgDelegate := stakingtypes.NewMsgDelegate(address, validator_address, amount)

	if err := msgMsgDelegate.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgDelegate")
	}

	_, err = msgServer.Delegate(ctx, msgMsgDelegate) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "elys stake msg")
	}

	var resp = &types.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Staking succeed",
	}

	return resp, nil
}

func performMsgCommit(f *commitmentkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgStake *types.MsgStake) (*types.RequestResponse, error) {
	if msgStake == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid staking parameter"}
	}
	msgServer := commitmentkeeper.NewMsgServerImpl(*f)
	msgMsgCommit := commitmenttypes.NewMsgCommitUnclaimedRewards(msgStake.Address, msgStake.Amount, msgStake.Asset)

	if err := msgMsgCommit.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgCommit")
	}

	_, err := msgServer.CommitUnclaimedRewards(ctx, msgMsgCommit) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "commit msg")
	}

	var resp = &types.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Staking succeed",
	}
	return resp, nil
}
