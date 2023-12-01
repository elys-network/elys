package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/stablestake/keeper"
	types "github.com/elys-network/elys/x/stablestake/types"
)

func (m *Messenger) msgUnbond(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgUnbond) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Unbond null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msg.Creator != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "unbond wrong sender"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msg")
	}

	res, err := msgServer.Unbond(
		sdk.WrapSDKContext(ctx),
		msg,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "Unbond msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize Unbond response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
