package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/commitment/keeper"
	types "github.com/elys-network/elys/x/commitment/types"
)

func (m *Messenger) msgWithdrawTokens(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgWithdrawTokens) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "WithdrawTokens null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgWithdrawTokens := types.NewMsgWithdrawTokens(
		msg.Creator,
		msg.Amount,
		msg.Denom,
	)

	if err := msgWithdrawTokens.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgWithdrawTokens")
	}

	res, err := msgServer.WithdrawTokens(
		sdk.WrapSDKContext(ctx),
		msgWithdrawTokens,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "WithdrawTokens msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize WithdrawTokens response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
