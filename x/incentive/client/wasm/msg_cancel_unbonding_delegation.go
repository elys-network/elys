package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (m *Messenger) msgCancelUnbondingDelegation(ctx sdk.Context, contractAddr sdk.AccAddress, msgCancelUnbonding *stakingtypes.MsgCancelUnbondingDelegation) ([]sdk.Event, [][]byte, error) {
	var res *wasmbindingstypes.RequestResponse
	var err error

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgCancelUnbonding.DelegatorAddress != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "wrong sender"}
	}

	if msgCancelUnbonding.Amount.Denom != paramtypes.Elys {
		return nil, nil, errorsmod.Wrap(err, "invalid asset!")
	}

	res, err = performMsgCancelUnbondingElys(m.stakingKeeper, ctx, contractAddr, msgCancelUnbonding)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform cancel elys unbonding")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize stake")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func performMsgCancelUnbondingElys(f *stakingkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgCancelUnbonding *stakingtypes.MsgCancelUnbondingDelegation) (*wasmbindingstypes.RequestResponse, error) {
	if msgCancelUnbonding == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid cancel unbonding parameter"}
	}

	msgServer := stakingkeeper.NewMsgServerImpl(f)
	address, err := sdk.AccAddressFromBech32(msgCancelUnbonding.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	valAddr, err := sdk.ValAddressFromBech32(msgCancelUnbonding.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	msgMsgCancelUnbonding := stakingtypes.NewMsgCancelUnbondingDelegation(address, valAddr, msgCancelUnbonding.CreationHeight, msgCancelUnbonding.Amount)

	if err := msgMsgCancelUnbonding.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgCancelUnbonding")
	}

	_, err = msgServer.CancelUnbondingDelegation(sdk.WrapSDKContext(ctx), msgMsgCancelUnbonding)
	if err != nil {
		return nil, errorsmod.Wrap(err, "elys cancel bonding msg")
	}

	var resp = &wasmbindingstypes.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Cancel unbonding succeed!",
	}

	return resp, nil
}
