package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	margintypes "github.com/elys-network/elys/x/margin/types"
)

func (m *Messenger) msgClose(ctx sdk.Context, contractAddr sdk.AccAddress, msgClose *margintypes.MsgClose) ([]sdk.Event, [][]byte, error) {
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

func PerformMsgClose(f *marginkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgClose *margintypes.MsgClose) (*margintypes.MsgCloseResponse, error) {
	if msgClose == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "margin close null margin close"}
	}
	msgServer := marginkeeper.NewMsgServerImpl(*f)

	msgMsgClose := margintypes.NewMsgClose(msgClose.Creator, uint64(msgClose.Id))

	if err := msgMsgClose.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgClose")
	}

	_, err := msgServer.Close(ctx, msgMsgClose) // Discard the response because it's empty

	if err != nil {
		return nil, errorsmod.Wrap(err, "margin close msg")
	}

	var resp = &margintypes.MsgCloseResponse{}
	return resp, nil
}
