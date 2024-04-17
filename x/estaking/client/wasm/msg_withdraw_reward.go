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

func (m *Messenger) msgWithdrawReward(ctx sdk.Context, contractAddr sdk.AccAddress, msgWithdrawReward *types.MsgWithdrawReward) ([]sdk.Event, [][]byte, error) {
	var res *wasmbindingstypes.RequestResponse
	var err error

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgWithdrawReward.DelegatorAddress != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "wrong sender"}
	}

	res, err = performMsgWithdrawReward(m.keeper, ctx, contractAddr, msgWithdrawReward)
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

func performMsgWithdrawReward(f *estakingkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgWithdrawReward *types.MsgWithdrawReward) (*wasmbindingstypes.RequestResponse, error) {
	if msgWithdrawReward == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid claim rewards parameter"}
	}

	msgServer := estakingkeeper.NewMsgServerImpl(*f)
	_, err := sdk.AccAddressFromBech32(msgWithdrawReward.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	msgMsgWithdrawReward := &types.MsgWithdrawReward{
		DelegatorAddress: msgWithdrawReward.DelegatorAddress,
		ValidatorAddress: msgWithdrawReward.ValidatorAddress,
	}

	if err := msgMsgWithdrawReward.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgWithdrawReward")
	}

	_, err = msgServer.WithdrawReward(sdk.WrapSDKContext(ctx), msgMsgWithdrawReward) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "elys withdrawReward msg")
	}

	resp := &wasmbindingstypes.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "WithdrawReward succeed!",
	}

	return resp, nil
}
