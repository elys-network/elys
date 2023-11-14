package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/commitment/keeper"
	types "github.com/elys-network/elys/x/commitment/types"
)

func (m *Messenger) msgCommitClaimedRewards(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgCommitClaimedRewards) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "CommitUnclaimedRewards null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgCommitClaimedRewards := types.NewMsgCommitClaimedRewards(
		msg.Creator,
		msg.Amount,
		msg.Denom,
	)

	if err := msgCommitClaimedRewards.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgCommitClaimedRewards")
	}

	res, err := msgServer.CommitClaimedRewards(
		sdk.WrapSDKContext(ctx),
		msgCommitClaimedRewards,
	)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "CommitUnclaimedRewards msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize CommitUnclaimedRewards response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
