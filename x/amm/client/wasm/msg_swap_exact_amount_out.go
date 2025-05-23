package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/v5/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/v5/x/amm/types"
)

func (m *Messenger) msgSwapExactAmountOut(ctx sdk.Context, contractAddr sdk.AccAddress, msg *ammtypes.MsgSwapExactAmountOut) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "swap exact amount out null msg"}
	}

	if msg.Sender != contractAddr.String() {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "swap exact amount out wrong sender"}
	}

	msgServer := ammkeeper.NewMsgServerImpl(*m.keeper)

	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating swap exact amount out msg")
	}

	res, err := msgServer.SwapExactAmountOut(ctx, msg)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "swap exact amount out msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize swap exact amount out response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
