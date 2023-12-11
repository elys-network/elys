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

func (m *Messenger) msgVest(ctx sdk.Context, contractAddr sdk.AccAddress, msgVest *commitmenttypes.MsgVest) ([]sdk.Event, [][]byte, error) {
	if msgVest == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Vest null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgVest.Creator != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "vest wrong sender"}
	}

	var res *wasmbindingstypes.RequestResponse
	var err error
	if msgVest.Denom != paramtypes.Eden {
		return nil, nil, errorsmod.Wrap(err, "invalid asset!")
	}

	res, err = performMsgVestEden(m.keeper, ctx, contractAddr, msgVest)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform eden vest")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize stake")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func performMsgVestEden(f *commitmentkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgVest *commitmenttypes.MsgVest) (*wasmbindingstypes.RequestResponse, error) {
	if msgVest == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid vesting parameter"}
	}

	msgServer := commitmentkeeper.NewMsgServerImpl(*f)
	msgMsgVest := commitmenttypes.NewMsgVest(msgVest.Creator, msgVest.Amount, msgVest.Denom)

	if err := msgMsgVest.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgVest")
	}

	_, err := msgServer.Vest(sdk.WrapSDKContext(ctx), msgMsgVest) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "eden vesting msg")
	}

	resp := &wasmbindingstypes.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Eden Vesting succeed!",
	}

	return resp, nil
}
