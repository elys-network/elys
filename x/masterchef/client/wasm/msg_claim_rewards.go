package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	masterchefkeeper "github.com/elys-network/elys/x/masterchef/keeper"
	"github.com/elys-network/elys/x/masterchef/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (m *Messenger) msgClaimRewards(ctx sdk.Context, contractAddr sdk.AccAddress, msgClaimRewards *types.MsgClaimRewards) ([]sdk.Event, [][]byte, error) {
	var res *wasmbindingstypes.RequestResponse
	var err error

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgClaimRewards.Sender != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "wrong sender"}
	}

	res, err = performMsgClaimRewards(m.keeper, ctx, contractAddr, msgClaimRewards)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform elys claim rewards")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize claim rewards response")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func performMsgClaimRewards(f *masterchefkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgClaimRewards *types.MsgClaimRewards) (*wasmbindingstypes.RequestResponse, error) {
	if msgClaimRewards == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid claim rewards parameter"}
	}

	msgServer := masterchefkeeper.NewMsgServerImpl(*f)
	_, err := sdk.AccAddressFromBech32(msgClaimRewards.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid claim rewards msg sender address")
	}

	msgMsgClaimRewards := &types.MsgClaimRewards{
		Sender:  msgClaimRewards.Sender,
		PoolIds: msgClaimRewards.PoolIds,
	}

	if err := msgMsgClaimRewards.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgClaimRewards")
	}

	_, err = msgServer.ClaimRewards(sdk.WrapSDKContext(ctx), msgMsgClaimRewards) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to process claim rewards msg")
	}

	resp := &wasmbindingstypes.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Claim rewards succeed!",
	}

	return resp, nil
}
