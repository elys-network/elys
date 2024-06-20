package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	perpetualkeeper "github.com/elys-network/elys/x/perpetual/keeper"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
)

func (m *Messenger) msgAddCollateral(ctx sdk.Context, contractAddr sdk.AccAddress, msgAddCollateral *perpetualtypes.MsgBrokerAddCollateral) ([]sdk.Event, [][]byte, error) {
	if msgAddCollateral == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Add collateral null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "contract address must be broker address"}
	}

	if msgAddCollateral.Creator != contractAddr.String() {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "add collateral wrong sender"}
	}

	msgServer := perpetualkeeper.NewMsgServerImpl(*m.keeper)

	if err := msgAddCollateral.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgMsgAddCollateral")
	}

	res, err := msgServer.BrokerAddCollateral(sdk.WrapSDKContext(ctx), msgAddCollateral)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perpetual add collateral msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize add collateral response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
