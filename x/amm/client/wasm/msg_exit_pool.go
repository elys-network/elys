package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func (m *Messenger) msgExitPool(ctx sdk.Context, contractAddr sdk.AccAddress, msg *ammtypes.MsgExitPool) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "exit pool null msg"}
	}

	msgServer := ammkeeper.NewMsgServerImpl(*m.keeper)

	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msg")
	}

	res, err := msgServer.ExitPool(
		sdk.WrapSDKContext(ctx),
		msg,
	)
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
