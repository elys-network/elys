package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/commitment/keeper"
	types "github.com/elys-network/elys/x/commitment/types"
)

func (m *Messenger) msgCommitUnclaimedRewards(ctx sdk.Context, contractAddr sdk.AccAddress, msg *types.MsgCommitUnclaimedRewards) ([]sdk.Event, [][]byte, error) {
	if msg == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "CommitUnclaimedRewards null msg"}
	}

	msgServer := keeper.NewMsgServerImpl(*m.keeper)

	msgCommitUnclaimedRewards := types.NewMsgCommitUnclaimedRewards(
		msg.Creator,
		msg.Amount,
		msg.Denom,
	)

	if err := msgCommitUnclaimedRewards.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgCommitUnclaimedRewards")
	}

	res, err := msgServer.CommitUnclaimedRewards(
		sdk.WrapSDKContext(ctx),
		msgCommitUnclaimedRewards,
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
