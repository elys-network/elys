package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	incentivekeeper "github.com/elys-network/elys/x/incentive/keeper"
	incentivetypes "github.com/elys-network/elys/x/incentive/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (m *Messenger) msgWithdrawValidatorCommission(ctx sdk.Context, contractAddr sdk.AccAddress, msgWithdrawValidatorCommission *incentivetypes.MsgWithdrawValidatorCommission) ([]sdk.Event, [][]byte, error) {
	var res *wasmbindingstypes.RequestResponse
	var err error

	brokerAddress := m.parameterKeeper.GetParams(ctx).BrokerAddress
	if msgWithdrawValidatorCommission.DelegatorAddress != contractAddr.String() && contractAddr.String() != brokerAddress {
		return nil, nil, wasmvmtypes.InvalidRequest{Err: "wrong sender"}
	}

	res, err = performMsgWithdrawValidatorCommissions(m.keeper, ctx, contractAddr, msgWithdrawValidatorCommission)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform withdraw validator commission")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize withdraw validator commission")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func performMsgWithdrawValidatorCommissions(f *incentivekeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgWithdrawValidatorCommission *incentivetypes.MsgWithdrawValidatorCommission) (*wasmbindingstypes.RequestResponse, error) {
	if msgWithdrawValidatorCommission == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid withdraw validator commission parameter"}
	}

	address, err := sdk.AccAddressFromBech32(msgWithdrawValidatorCommission.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	valAddr, err := sdk.ValAddressFromBech32(msgWithdrawValidatorCommission.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	msgServer := incentivekeeper.NewMsgServerImpl(*f)
	msgMsgWithdrawValidatorCommissions := incentivetypes.NewMsgWithdrawValidatorCommission(address, valAddr)

	if err := msgMsgWithdrawValidatorCommissions.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgWithdrawValidatorCommission")
	}

	_, err = msgServer.WithdrawValidatorCommission(sdk.WrapSDKContext(ctx), msgMsgWithdrawValidatorCommissions) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "withdraw validator commission msg")
	}

	resp := &wasmbindingstypes.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Withdraw validator commissions succeed!",
	}

	return resp, nil
}
