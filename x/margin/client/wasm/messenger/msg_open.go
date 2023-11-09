package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	cosmos_sdk_math "cosmossdk.io/math"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/client/wasm/types"
	marginkeeper "github.com/elys-network/elys/x/margin/keeper"
	margintypes "github.com/elys-network/elys/x/margin/types"
)

func (m *Messenger) msgOpen(ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *types.MsgOpen) ([]sdk.Event, [][]byte, error) {
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

func PerformMsgOpen(f *marginkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgOpen *types.MsgOpen) (*types.MsgOpenResponse, error) {
	if msgOpen == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "margin open null margin open"}
	}
	msgServer := marginkeeper.NewMsgServerImpl(*f)

	msgMsgOpen := margintypes.NewMsgOpen(msgOpen.Creator, msgOpen.CollateralAsset, cosmos_sdk_math.Int(msgOpen.CollateralAmount), msgOpen.BorrowAsset, msgOpen.Position, msgOpen.Leverage, msgOpen.TakeProfitPrice)

	if err := msgMsgOpen.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgOpen")
	}

	_, err := msgServer.Open(ctx, msgMsgOpen) // Discard the response because it's empty

	if err != nil {
		return nil, errorsmod.Wrap(err, "margin open msg")
	}

	var resp = &types.MsgOpenResponse{
		MetaData: msgOpen.MetaData,
	}
	return resp, nil
}
