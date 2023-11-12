package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func (m *Messenger) msgCreatePool(ctx sdk.Context, contractAddr sdk.AccAddress, msg *ammtypes.MsgCreatePool) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "create pool null msg"}
	}

	msgServer := ammkeeper.NewMsgServerImpl(*m.keeper)

	msgCreatePool := ammtypes.NewMsgCreatePool(msg.Sender, msg.PoolParams, msg.PoolAssets)

	if err := msgCreatePool.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgCreatePool")
	}

	res, err := msgServer.CreatePool(
		sdk.WrapSDKContext(ctx),
		msgCreatePool,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "create pool msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize create pool response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
