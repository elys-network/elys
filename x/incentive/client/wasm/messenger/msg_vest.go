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

func (m *Messenger) msgVest(ctx sdk.Context, contractAddr sdk.AccAddress, msgVest *types.MsgVest) ([]sdk.Event, [][]byte, error) {
	var res *types.RequestResponse
	var err error
	if msgVest.Denom != paramtypes.Eden {
		return nil, nil, errorsmod.Wrap(err, "invalid asset!")
	}

	res, err = performMsgVestEden(m.commitmentKeeper, ctx, contractAddr, msgVest)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform eden vest")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize stake")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func performMsgVestEden(f *commitmentkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgVest *types.MsgVest) (*types.RequestResponse, error) {
	if msgVest == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid vesting parameter"}
	}

	msgServer := commitmentkeeper.NewMsgServerImpl(*f)
	msgMsgVest := commitmenttypes.NewMsgVest(msgVest.Creator, msgVest.Amount, msgVest.Denom)

	if err := msgMsgVest.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgVest")
	}

	_, err := msgServer.Vest(ctx, msgMsgVest) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "eden vesting msg")
	}

	var resp = &types.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Eden Vesting succeed!",
	}

	return resp, nil
}
