package messenger

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/client/wasm/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (m *Messenger) msgCancelVest(ctx sdk.Context, contractAddr sdk.AccAddress, msgCancelVest *types.MsgCancelVest) ([]sdk.Event, [][]byte, error) {
	var res *types.RequestResponse
	var err error
	if msgCancelVest.Denom != paramtypes.Eden {
		return nil, nil, errorsmod.Wrap(err, "invalid asset!")
	}

	res, err = performMsgCancelVestEden(m.commitmentKeeper, ctx, contractAddr, msgCancelVest)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform eden cancel vest")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize cancel vesting")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func performMsgCancelVestEden(f *commitmentkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgCancelVest *types.MsgCancelVest) (*types.RequestResponse, error) {
	if msgCancelVest == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid cancel vesting parameter"}
	}

	msgServer := commitmentkeeper.NewMsgServerImpl(*f)
	msgMsgCancelVest := commitmenttypes.NewMsgCancelVest(msgCancelVest.Creator, msgCancelVest.Amount, msgCancelVest.Denom)

	if err := msgMsgCancelVest.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgCancelVest")
	}

	_, err := msgServer.CancelVest(ctx, msgMsgCancelVest) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "eden vesting msg")
	}

	var resp = &types.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Eden vesting cancel succeed!",
	}

	return resp, nil
}
