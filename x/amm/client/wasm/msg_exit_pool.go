package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/v7/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
)

func (m *Messenger) msgExitPool(ctx sdk.Context, contractAddr sdk.AccAddress, msg *ammtypes.MsgExitPool) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "exit pool null msg"}
	}

	if msg.Sender != contractAddr.String() {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "exit pool wrong sender"}
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msg")
	}

	msgServer := ammkeeper.NewMsgServerImpl(*m.keeper)
	res, err := msgServer.ExitPool(ctx, msg)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "exit pool msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize exit pool response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
