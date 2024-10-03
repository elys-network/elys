package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	leveragelpkeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
)

func (m *Messenger) msgClose(ctx sdk.Context, contractAddr sdk.AccAddress, msgClose *leveragelptypes.MsgClose) ([]sdk.Event, [][]byte, error) {
	if msgClose == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Close null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgClose.Creator != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "close wrong sender"}
	}

	if err := msgClose.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgMsgClose")
	}

	msgServer := leveragelpkeeper.NewMsgServerImpl(*m.keeper)

	res, err := msgServer.Close(sdk.WrapSDKContext(ctx), msgClose)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "leveragelp close msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize close response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
