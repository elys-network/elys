package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func (m *Messenger) msgSwapExactAmountIn(ctx sdk.Context, contractAddr sdk.AccAddress, msgSwapExactAmountIn *ammtypes.MsgSwapExactAmountIn) ([]sdk.Event, [][]byte, error) {
	res, err := performMsgSwapExactAmountIn(m.keeper, ctx, contractAddr, msgSwapExactAmountIn)
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

func performMsgSwapExactAmountIn(f *ammkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgSwapExactAmountIn *ammtypes.MsgSwapExactAmountIn) (*ammtypes.MsgSwapExactAmountInResponse, error) {
	if msgSwapExactAmountIn == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "swap null swap"}
	}

	msgServer := ammkeeper.NewMsgServerImpl(*f)

	var PoolIds []uint64
	var TokenOutDenoms []string

	for _, route := range msgSwapExactAmountIn.Routes {
		PoolIds = append(PoolIds, route.PoolId)
		TokenOutDenoms = append(TokenOutDenoms, route.TokenOutDenom)
	}

	msgMsgSwapExactAmountIn := ammtypes.NewMsgSwapExactAmountIn(msgSwapExactAmountIn.Sender, msgSwapExactAmountIn.TokenIn, msgSwapExactAmountIn.TokenOutMinAmount, PoolIds, TokenOutDenoms)

	if err := msgMsgSwapExactAmountIn.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating MsgMsgSwapExactAmountIn")
	}

	// Swap
	res, err := msgServer.SwapExactAmountIn(
		sdk.WrapSDKContext(ctx),
		msgMsgSwapExactAmountIn,
	)
	if err != nil {
		return nil, errorsmod.Wrap(err, "swap msg")
	}

	return res, nil
}
