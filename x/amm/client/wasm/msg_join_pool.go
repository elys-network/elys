package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func (m *Messenger) msgJoinPool(ctx sdk.Context, contractAddr sdk.AccAddress, msg *ammtypes.MsgJoinPool) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "join pool null msg"}
	}

	msgServer := ammkeeper.NewMsgServerImpl(*m.keeper)

	msgJoinPool := ammtypes.NewMsgJoinPool(msg.Sender, msg.PoolId, msg.MaxAmountsIn, msg.ShareAmountOut, msg.NoRemaining)

	if err := msgJoinPool.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgJoinPool")
	}

	res, err := msgServer.JoinPool(
		sdk.WrapSDKContext(ctx),
		msgJoinPool,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "join pool msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize join pool response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
