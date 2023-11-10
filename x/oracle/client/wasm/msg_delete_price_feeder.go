package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/oracle/keeper"
	types "github.com/elys-network/elys/x/oracle/types"
)

func (m *Messenger) msgDeletePriceFeeder(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgDeletePriceFeeder) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "DeletePriceFeeder null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgDeletePriceFeeder := types.NewMsgDeletePriceFeeder(
		msg.Feeder,
	)

	if err := msgDeletePriceFeeder.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgDeletePriceFeeder")
	}

	res, err := msgServer.DeletePriceFeeder(
		sdk.WrapSDKContext(ctx),
		msgDeletePriceFeeder,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "DeletePriceFeeder msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize DeletePriceFeeder response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
