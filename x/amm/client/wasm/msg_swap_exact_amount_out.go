package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func (m *Messenger) msgSwapExactAmountOut(ctx sdk.Context, contractAddr sdk.AccAddress, msgSwapExactAmountOut *ammtypes.MsgSwapExactAmountOut) ([]sdk.Event, [][]byte, error) {
	res, err := performMsgSwapExactAmountOut(m.keeper, ctx, contractAddr, msgSwapExactAmountOut)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform swap")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize swap response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func performMsgSwapExactAmountOut(f *ammkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgSwapExactAmountOut *ammtypes.MsgSwapExactAmountOut) (*ammtypes.MsgSwapExactAmountOutResponse, error) {
	if msgSwapExactAmountOut == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "swap null swap"}
	}

	msgServer := ammkeeper.NewMsgServerImpl(*f)

	var PoolIds []uint64
	var TokenInDenoms []string

	for _, route := range msgSwapExactAmountOut.Routes {
		PoolIds = append(PoolIds, route.PoolId)
		TokenInDenoms = append(TokenInDenoms, route.TokenInDenom)
	}

	msgMsgSwapExactAmountOut := ammtypes.NewMsgSwapExactAmountOut(msgSwapExactAmountOut.Sender, msgSwapExactAmountOut.TokenOut, msgSwapExactAmountOut.TokenInMaxAmount, PoolIds, TokenInDenoms)

	if err := msgMsgSwapExactAmountOut.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating MsgMsgSwapExactAmountOut")
	}

	// Swap
	res, err := msgServer.SwapExactAmountOut(
		sdk.WrapSDKContext(ctx),
		msgMsgSwapExactAmountOut,
	)
	if err != nil {
		return nil, errorsmod.Wrap(err, "swap msg")
	}

	return res, nil
}
