package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	estakingkeeper "github.com/elys-network/elys/x/estaking/keeper"
	"github.com/elys-network/elys/x/estaking/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (m *Messenger) msgWithdrawElysStakingRewards(ctx sdk.Context, contractAddr sdk.AccAddress, msgWithdrawElysStakingRewards *types.MsgWithdrawElysStakingRewards) ([]sdk.Event, [][]byte, error) {
	var res *wasmbindingstypes.RequestResponse
	var err error

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgWithdrawElysStakingRewards.DelegatorAddress != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "wrong sender"}
	}

	res, err = performMsgWithdrawElysStakingRewards(m.keeper, ctx, contractAddr, msgWithdrawElysStakingRewards)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform elys claim rewards")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize stake")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func performMsgWithdrawElysStakingRewards(f *estakingkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgWithdrawElysStakingRewards *types.MsgWithdrawElysStakingRewards) (*wasmbindingstypes.RequestResponse, error) {
	if msgWithdrawElysStakingRewards == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid claim elys staking rewards parameter"}
	}

	msgServer := estakingkeeper.NewMsgServerImpl(*f)
	_, err := sdk.AccAddressFromBech32(msgWithdrawElysStakingRewards.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	msg := &types.MsgWithdrawElysStakingRewards{
		DelegatorAddress: msgWithdrawElysStakingRewards.DelegatorAddress,
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating MsgWithdrawElysStakingRewards")
	}

	_, err = msgServer.WithdrawElysStakingRewards(sdk.WrapSDKContext(ctx), msg) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "elys WithdrawElysStakingRewards msg")
	}

	resp := &wasmbindingstypes.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "WithdrawElysStakingRewards succeed!",
	}

	return resp, nil
}
