package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	margintypes "github.com/elys-network/elys/x/margin/types"
)

func (m *Messenger) msgBrokerClose(ctx sdk.Context, contractAddr sdk.AccAddress, msgBrokerClose *margintypes.MsgBrokerClose) ([]sdk.Event, [][]byte, error) {
	if msgBrokerClose == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Broker Close null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgBrokerClose.Creator != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "broker close wrong sender"}
	}

	res, err := PerformMsgBrokerClose(m.keeper, ctx, contractAddr, msgBrokerClose)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform broker close")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize broker close response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgBrokerClose(f *marginkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgBrokerClose *margintypes.MsgBrokerClose) (*margintypes.MsgBrokerCloseResponse, error) {
	if msgBrokerClose == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "margin broker close null margin broker close"}
	}
	msgServer := marginkeeper.NewMsgServerImpl(*f)

	msgMsgBrokerClose := margintypes.NewMsgBrokerClose(msgBrokerClose.Creator, uint64(msgBrokerClose.Id), msgBrokerClose.Owner, msgBrokerClose.Amount)

	if err := msgMsgBrokerClose.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgBrokerClose")
	}

	_, err := msgServer.BrokerClose(sdk.WrapSDKContext(ctx), msgMsgBrokerClose) // Discard the response because it's empty

	if err != nil {
		return nil, errorsmod.Wrap(err, "margin broker close msg")
	}

	var resp = &margintypes.MsgBrokerCloseResponse{}
	return resp, nil
}
