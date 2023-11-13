package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/commitment/keeper"
	types "github.com/elys-network/elys/x/commitment/types"
)

func (m *Messenger) msgClaimReward(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgClaimReward) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "ClaimReward null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgClaimReward := types.NewMsgClaimReward(
		msg.Creator,
		msg.Amount,
		msg.Denom,
	)

	if err := msgClaimReward.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgClaimReward")
	}

	res, err := msgServer.ClaimReward(
		sdk.WrapSDKContext(ctx),
		msgClaimReward,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "ClaimReward msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize ClaimReward response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
