package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	margintypes "github.com/elys-network/elys/x/margin/types"
)

func (m *Messenger) msgOpen(ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *margintypes.MsgOpen) ([]sdk.Event, [][]byte, error) {
	if msgOpen == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Open null msg"}
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

func PerformMsgOpen(f *marginkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *margintypes.MsgOpen) (*margintypes.MsgOpenResponse, error) {
	if msgOpen == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "margin open null margin open"}
	}
	msgServer := marginkeeper.NewMsgServerImpl(*f)

	msgMsgOpen := margintypes.NewMsgOpen(
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

	_, err := msgServer.Open(ctx, msgMsgOpen) // Discard the response because it's empty

	if err != nil {
		return nil, errorsmod.Wrap(err, "margin open msg")
	}

	var resp = &margintypes.MsgOpenResponse{}
	return resp, nil
}
