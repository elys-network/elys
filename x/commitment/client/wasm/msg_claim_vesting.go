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

func (m *Messenger) msgClaimVesting(ctx sdk.Context, contractAddr sdk.AccAddress, msgClaimVesting *commitmenttypes.MsgClaimVesting) ([]sdk.Event, [][]byte, error) {
	if msgClaimVesting == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "ClaimVesting null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgClaimVesting.Sender != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "vest wrong sender"}
	}

	var res *wasmbindingstypes.RequestResponse
	var err error

	res, err = performMsgClaimVesting(m.keeper, ctx, contractAddr, msgClaimVesting)
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

func performMsgClaimVesting(f *commitmentkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgClaimVesting *commitmenttypes.MsgClaimVesting) (*wasmbindingstypes.RequestResponse, error) {
	if msgClaimVesting == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid vesting parameter"}
	}

	msgServer := commitmentkeeper.NewMsgServerImpl(*f)
	msgMsgClaimVesting := commitmenttypes.NewMsgClaimVesting(msgClaimVesting.Sender)

	if err := msgMsgClaimVesting.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgClaimVesting")
	}

	_, err := msgServer.ClaimVesting(sdk.WrapSDKContext(ctx), msgMsgClaimVesting) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "eden vesting msg")
	}

	resp := &wasmbindingstypes.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Eden Vesting succeed!",
	}

	return resp, nil
}
