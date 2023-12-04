package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (m *Messenger) msgCancelVest(ctx sdk.Context, contractAddr sdk.AccAddress, msgCancelVest *commitmenttypes.MsgCancelVest) ([]sdk.Event, [][]byte, error) {
	if msgCancelVest == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "cancel vest null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgCancelVest.Creator != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "cancel vest wrong sender"}
	}

	var res *wasmbindingstypes.RequestResponse
	var err error
	if msgCancelVest.Denom != paramtypes.Eden {
		return nil, nil, errorsmod.Wrap(err, "invalid asset!")
	}

	res, err = performMsgCancelVestEden(m.keeper, ctx, contractAddr, msgCancelVest)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform eden cancel vest")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize cancel vesting")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func performMsgCancelVestEden(f *commitmentkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgCancelVest *commitmenttypes.MsgCancelVest) (*wasmbindingstypes.RequestResponse, error) {
	if msgCancelVest == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid cancel vesting parameter"}
	}

	msgServer := commitmentkeeper.NewMsgServerImpl(*f)
	msgMsgCancelVest := commitmenttypes.NewMsgCancelVest(msgCancelVest.Creator, msgCancelVest.Amount, msgCancelVest.Denom)

	if err := msgMsgCancelVest.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgCancelVest")
	}

	_, err := msgServer.CancelVest(sdk.WrapSDKContext(ctx), msgMsgCancelVest) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "eden vesting msg")
	}

	var resp = &wasmbindingstypes.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Eden vesting cancel succeed!",
	}

	return resp, nil
}
