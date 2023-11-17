package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	cosmos_sdk_math "cosmossdk.io/math"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	margintypes "github.com/elys-network/elys/x/margin/types"
)

func (m *Messenger) msgBrokerOpen(ctx sdk.Context, contractAddr sdk.AccAddress, msgBrokerOpen *margintypes.MsgBrokerOpen) ([]sdk.Event, [][]byte, error) {
	res, err := PerformMsgBrokerOpen(m.keeper, ctx, contractAddr, msgBrokerOpen)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform broker open")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize broker open response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgBrokerOpen(f *marginkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgBrokerOpen *margintypes.MsgBrokerOpen) (*margintypes.MsgBrokerOpenResponse, error) {
	if msgBrokerOpen == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "margin broker open null margin broker open"}
	}
	msgServer := marginkeeper.NewMsgServerImpl(*f)

	msgMsgBrokerOpen := margintypes.NewMsgBrokerOpen(msgBrokerOpen.Creator, msgBrokerOpen.CollateralAsset, cosmos_sdk_math.Int(msgBrokerOpen.CollateralAmount), msgBrokerOpen.BorrowAsset, msgBrokerOpen.Position, msgBrokerOpen.Leverage, msgBrokerOpen.TakeProfitPrice, msgBrokerOpen.Owner)

	if err := msgMsgBrokerOpen.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgBrokerOpen")
	}

	_, err := msgServer.BrokerOpen(ctx, msgMsgBrokerOpen) // Discard the response because it's empty

	if err != nil {
		return nil, errorsmod.Wrap(err, "margin broker open msg")
	}

	var resp = &margintypes.MsgBrokerOpenResponse{}
	return resp, nil
}
