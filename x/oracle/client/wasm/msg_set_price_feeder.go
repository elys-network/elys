package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/oracle/keeper"
	types "github.com/elys-network/elys/x/oracle/types"
)

func (m *Messenger) msgSetPriceFeeder(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgSetPriceFeeder) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "SetPriceFeeder null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgSetPriceFeeder := types.NewMsgSetPriceFeeder(
		msg.Feeder,
		msg.IsActive,
	)

	if err := msgSetPriceFeeder.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgSetPriceFeeder")
	}

	res, err := msgServer.SetPriceFeeder(
		sdk.WrapSDKContext(ctx),
		msgSetPriceFeeder,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "SetPriceFeeder msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize SetPriceFeeder response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
