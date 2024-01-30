package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	perpetualkeeper "github.com/elys-network/elys/x/perpetual/keeper"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
)

func (m *Messenger) msgOpen(ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *perpetualtypes.MsgOpen) ([]sdk.Event, [][]byte, error) {
	if msgOpen == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Open null msg"}
	}

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgOpen.Creator != contractAddr.String() && contractAddr.String() != brokerAddress {
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

func PerformMsgOpen(f *perpetualkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *perpetualtypes.MsgOpen) (*perpetualtypes.MsgOpenResponse, error) {
	if msgOpen == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "perpetual open null perpetual open"}
	}
	msgServer := perpetualkeeper.NewMsgServerImpl(*f)

	msgMsgOpen := perpetualtypes.NewMsgOpen(
		msgOpen.Creator,
		msgOpen.Position,
		msgOpen.Leverage,
		msgOpen.TradingAsset,
		msgOpen.Collateral,
		msgOpen.TakeProfitPrice,
	)

	if err := msgMsgOpen.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgOpen")
	}

	res, err := msgServer.Open(sdk.WrapSDKContext(ctx), msgMsgOpen) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "perpetual open msg")
	}

	resp := &perpetualtypes.MsgOpenResponse{
		Id: res.Id,
	}
	return resp, nil
}
