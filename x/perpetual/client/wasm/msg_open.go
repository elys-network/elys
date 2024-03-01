package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
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

	res, err := PerformMsgOpen(m.keeper, ctx, contractAddr, msgOpen)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform open")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize open response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func PerformMsgOpen(f *perpetualkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *perpetualtypes.MsgBrokerOpen) (*perpetualtypes.MsgOpenResponse, error) {
	if msgOpen == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "perpetual open null perpetual open"}
	}
	msgServer := perpetualkeeper.NewMsgServerImpl(*f)

	msgMsgOpen := perpetualtypes.NewMsgBrokerOpen(
		msgOpen.Creator,
		msgOpen.Position,
		msgOpen.Leverage,
		msgOpen.TradingAsset,
		msgOpen.Collateral,
		msgOpen.TakeProfitPrice,
		msgOpen.Owner,
	)

	if err := msgMsgOpen.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgOpen")
	}

	res, err := msgServer.BrokerOpen(sdk.WrapSDKContext(ctx), msgMsgOpen) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "perpetual open msg")
	}

	resp := &perpetualtypes.MsgOpenResponse{
		Id: res.Id,
	}
	return resp, nil
}
