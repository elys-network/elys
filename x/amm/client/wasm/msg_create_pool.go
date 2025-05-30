package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/v6/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
)

func (m *Messenger) msgCreatePool(ctx sdk.Context, contractAddr sdk.AccAddress, msg *ammtypes.MsgCreatePool) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "create pool null msg"}
	}

	if msg.Sender != contractAddr.String() {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "create pool wrong sender"}
	}

	msgServer := ammkeeper.NewMsgServerImpl(*m.keeper)

	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msg")
	}

	res, err := msgServer.CreatePool(ctx, msg)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "create pool msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize create pool response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
