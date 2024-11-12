package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	perpetualkeeper "github.com/elys-network/elys/x/perpetual/keeper"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
)

func (m *Messenger) msgOpen(ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *perpetualtypes.MsgBrokerOpen) ([]sdk.Event, [][]byte, error) {
	if msgOpen == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Open null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "contract address must be broker address"}
	}

	if msgOpen.Creator != contractAddr.String() {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "open wrong sender"}
	}

	msgServer := perpetualkeeper.NewMsgServerImpl(*m.keeper)

	if err := msgOpen.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgMsgOpen")
	}

	res, err := msgServer.BrokerOpen(sdk.WrapSDKContext(ctx), msgOpen)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perpetual open msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize open response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
