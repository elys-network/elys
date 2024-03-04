package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	perpetualkeeper "github.com/elys-network/elys/x/perpetual/keeper"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
)

func (m *Messenger) msgClose(ctx sdk.Context, contractAddr sdk.AccAddress, msgClose *perpetualtypes.MsgBrokerClose) ([]sdk.Event, [][]byte, error) {
	if msgClose == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Close null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "contract address must be broker address"}
	}

	if msgClose.Creator != contractAddr.String() {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "close wrong sender"}
	}

	msgServer := perpetualkeeper.NewMsgServerImpl(*m.keeper)

	if err := msgClose.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgMsgClose")
	}

	res, err := msgServer.BrokerClose(sdk.WrapSDKContext(ctx), msgClose)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perpetual close msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize close response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
