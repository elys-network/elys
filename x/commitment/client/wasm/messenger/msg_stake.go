package messenger

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
)

func (m *Messenger) msgStake(ctx sdk.Context, contractAddr sdk.AccAddress, msgStake *commitmenttypes.MsgStake) ([]sdk.Event, [][]byte, error) {
	if msgStake == nil {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "Invalid staking parameter"}
	}

	msgServer := commitmentkeeper.NewMsgServerImpl(*m.keeper)
	msgMsgStake := commitmenttypes.NewMsgStake(msgStake.Creator, msgStake.Amount, msgStake.Asset, msgStake.ValidatorAddress)

	if err := msgMsgStake.ValidateBasic(); err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed validating msgMsgStake")
	}

	res, err := msgServer.Stake(ctx, msgMsgStake)
	if err != nil { // Discard the response because it's empty
		return nil, nil, errorsmod.Wrap(err, "elys staking msg")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize stake")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}
