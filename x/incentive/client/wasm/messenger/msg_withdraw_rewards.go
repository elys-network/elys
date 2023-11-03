package messenger

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/incentive/client/wasm/types"
	incentivekeeper "github.com/elys-network/elys/x/incentive/keeper"
	incentivetypes "github.com/elys-network/elys/x/incentive/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (m *Messenger) msgWithdrawRewards(ctx sdk.Context, contractAddr sdk.AccAddress, msgWithdrawRewards *types.MsgWithdrawRewards) ([]sdk.Event, [][]byte, error) {
	var res *types.RequestResponse
	var err error

	res, err = performMsgWithdrawRewards(m.incentiveKeeper, ctx, contractAddr, msgWithdrawRewards)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform withdraw rewards")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize withdraw rewards")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func performMsgWithdrawRewards(f *incentivekeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgWithdrawRewards *types.MsgWithdrawRewards) (*types.RequestResponse, error) {
	if msgWithdrawRewards == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid withdraw rewards parameter"}
	}

	address, err := sdk.AccAddressFromBech32(msgWithdrawRewards.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	msgServer := incentivekeeper.NewMsgServerImpl(*f)
	msgMsgWithdrawRewards := incentivetypes.NewMsgWithdrawRewards(address, msgWithdrawRewards.Denom)

	if err := msgMsgWithdrawRewards.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgWithdrawRewards")
	}

	_, err = msgServer.WithdrawRewards(ctx, msgMsgWithdrawRewards) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "withdraw rewards msg")
	}

	var resp = &types.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Withdraw rewards succeed!",
	}

	return resp, nil
}
