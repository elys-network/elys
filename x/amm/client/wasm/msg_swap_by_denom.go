package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/v7/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
)

func (m *Messenger) msgSwapByDenom(ctx sdk.Context, contractAddr sdk.AccAddress, msg *ammtypes.MsgSwapByDenom) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "swap by denom null msg"}
	}

	if msg.Sender != contractAddr.String() {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "swap by denom wrong sender"}
	}

	msgServer := ammkeeper.NewMsgServerImpl(*m.keeper)

	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating swap by denom msg")
	}

	res, err := msgServer.SwapByDenom(ctx, msg)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "swap by denom msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize swap by denom response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
