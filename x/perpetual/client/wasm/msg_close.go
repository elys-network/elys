package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	perpetualkeeper "github.com/elys-network/elys/x/perpetual/keeper"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
)

func (m *Messenger) msgClose(ctx sdk.Context, contractAddr sdk.AccAddress, msgClose *perpetualtypes.MsgClose) ([]sdk.Event, [][]byte, error) {
	if msgClose == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Close null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgClose.Creator != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "close wrong sender"}
	}

	res, err := PerformMsgClose(m.keeper, ctx, contractAddr, msgClose)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform close")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize close response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgClose(f *perpetualkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgClose *perpetualtypes.MsgClose) (*perpetualtypes.MsgCloseResponse, error) {
	if msgClose == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "perpetual close null perpetual close"}
	}
	msgServer := perpetualkeeper.NewMsgServerImpl(*f)

	msgMsgClose := perpetualtypes.NewMsgClose(msgClose.Creator, uint64(msgClose.Id), msgClose.Amount)

	if err := msgMsgClose.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgClose")
	}

	_, err := msgServer.Close(sdk.WrapSDKContext(ctx), msgMsgClose) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "perpetual close msg")
	}

	resp := &perpetualtypes.MsgCloseResponse{}
	return resp, nil
}
