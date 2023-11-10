package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/commitment/keeper"
	types "github.com/elys-network/elys/x/commitment/types"
)

func (m *Messenger) msgCommitLiquidTokens(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgCommitLiquidTokens) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "CommitLiquidTokens null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgCommitLiquidTokens := types.NewMsgCommitLiquidTokens(
		msg.Creator,
		msg.Amount,
		msg.Denom,
	)

	if err := msgCommitLiquidTokens.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgCommitLiquidTokens")
	}

	res, err := msgServer.CommitLiquidTokens(
		sdk.WrapSDKContext(ctx),
		msgCommitLiquidTokens,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "CommitLiquidTokens msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize CommitLiquidTokens response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
