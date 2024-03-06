package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func (m *Messenger) msgJoinPool(ctx sdk.Context, contractAddr sdk.AccAddress, msg *ammtypes.MsgJoinPool) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "join pool null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msg.Sender != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "join pool wrong sender"}
	}

	msgServer := ammkeeper.NewMsgServerImpl(*m.keeper)

	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msg")
	}

	res, err := msgServer.JoinPool(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "join pool msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize join pool response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
