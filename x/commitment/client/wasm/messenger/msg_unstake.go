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

func (m *Messenger) msgUnstake(ctx sdk.Context, contractAddr sdk.AccAddress, msgUnstake *types.MsgUnstake) ([]sdk.Event, [][]byte, error) {
	var res *types.RequestResponse
	var err error
	if msgUnstake.Asset == paramtypes.Elys {
		res, err = performMsgUnstakeElys(m.stakingKeeper, ctx, contractAddr, msgUnstake)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "perform elys stake")
		}
	} else {
		res, err = performMsgUncommit(m.keeper, ctx, contractAddr, msgUnstake)
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

func performMsgUnstakeElys(f *stakingkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgUnstake *types.MsgUnstake) (*types.RequestResponse, error) {
	if msgUnstake == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid unstaking parameter"}
	}

	msgServer := stakingkeeper.NewMsgServerImpl(f)
	address, err := sdk.AccAddressFromBech32(msgUnstake.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	validator_address, err := sdk.ValAddressFromBech32(msgUnstake.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	amount := sdk.NewCoin(msgUnstake.Asset, msgUnstake.Amount)
	msgMsgUndelegate := stakingtypes.NewMsgUndelegate(address, validator_address, amount)

	if err := msgMsgUndelegate.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgDelegate")
	}

	_, err = msgServer.Undelegate(ctx, msgMsgUndelegate) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "elys unstake msg")
	}

	var resp = &types.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Unstaking succeed",
	}

	return resp, nil
}

func performMsgUncommit(f *commitmentkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgUnstake *types.MsgUnstake) (*types.RequestResponse, error) {
	if msgUnstake == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid staking parameter"}
	}
	msgServer := commitmentkeeper.NewMsgServerImpl(*f)
	msgMsgUncommit := commitmenttypes.NewMsgUncommitTokens(msgUnstake.Address, msgUnstake.Amount, msgUnstake.Asset)

	if err := msgMsgUncommit.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgCommit")
	}

	_, err := msgServer.UncommitTokens(ctx, msgMsgUncommit) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "commit msg")
	}

	var resp = &types.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Unstaking succeed",
	}
	return resp, nil
}
