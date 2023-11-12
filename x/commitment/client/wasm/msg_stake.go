package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (m *Messenger) msgStake(ctx sdk.Context, contractAddr sdk.AccAddress, msgStake *commitmenttypes.MsgStake) ([]sdk.Event, [][]byte, error) {
	if msgStake == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Invalid staking parameter"}
	}

	msgServer := commitmentkeeper.NewMsgServerImpl(*m.keeper)
	msgMsgStake := commitmenttypes.NewMsgStake(msgStake.Creator, msgStake.Amount, msgStake.Asset, msgStake.ValidatorAddress)

	if err := msgMsgStake.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgMsgStake")
	}

	res, err := msgServer.Stake(ctx, msgMsgStake)
	if err != nil { // Discard the response because it's empty
		return nil, nil, errorsmod.Wrap(err, "elys staking msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize stake")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func performMsgStakeElys(f *stakingkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgStake *wasmbindingstypes.MsgStake) (*wasmbindingstypes.RequestResponse, error) {
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

	var resp = &wasmbindingstypes.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Staking succeed",
	}

	return resp, nil
}

func performMsgCommit(f *commitmentkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgStake *wasmbindingstypes.MsgStake) (*wasmbindingstypes.RequestResponse, error) {
	if msgStake == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid staking parameter"}
	}
	msgServer := commitmentkeeper.NewMsgServerImpl(*f)
	msgMsgCommit := commitmenttypes.NewMsgCommitClaimedRewards(msgStake.Address, msgStake.Amount, msgStake.Asset)

	if err := msgMsgCommit.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgCommit")
	}

	_, err := msgServer.CommitClaimedRewards(ctx, msgMsgCommit) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "commit msg")
	}

	var resp = &wasmbindingstypes.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Staking succeed",
	}
	return resp, nil
}
