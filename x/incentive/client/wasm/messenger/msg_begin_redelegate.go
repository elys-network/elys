package messenger

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

func (m *Messenger) msgBeginRedelegate(ctx sdk.Context, contractAddr sdk.AccAddress, msgRedelegate *stakingtypes.MsgBeginRedelegate) ([]sdk.Event, [][]byte, error) {
	var res *wasmbindingstypes.RequestResponse
	var err error
	if msgRedelegate.Amount.Denom != paramtypes.Elys {
		return nil, nil, errorsmod.Wrap(err, "invalid asset!")
	}

	res, err = performMsgRedelegateElys(m.stakingKeeper, ctx, contractAddr, msgRedelegate)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform elys redelegate")
	}

	responseBytes, err := json.Marshal(*res)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to serialize stake")
	}

	resp := [][]byte{responseBytes}

	return nil, resp, nil
}

func performMsgRedelegateElys(f *stakingkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msgRedelegate *stakingtypes.MsgBeginRedelegate) (*wasmbindingstypes.RequestResponse, error) {
	if msgRedelegate == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "Invalid redelegate parameter"}
	}

	msgServer := stakingkeeper.NewMsgServerImpl(f)
	address, err := sdk.AccAddressFromBech32(msgRedelegate.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	valSrcAddr, err := sdk.ValAddressFromBech32(msgRedelegate.ValidatorSrcAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	valDstAddr, err := sdk.ValAddressFromBech32(msgRedelegate.ValidatorDstAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid address")
	}

	msgMsgRedelegate := stakingtypes.NewMsgBeginRedelegate(address, valSrcAddr, valDstAddr, msgRedelegate.Amount)

	if err := msgMsgRedelegate.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating msgMsgDelegate")
	}

	_, err = msgServer.BeginRedelegate(ctx, msgMsgRedelegate) // Discard the response because it's empty
	if err != nil {
		return nil, errorsmod.Wrap(err, "elys redelegation msg")
	}

	var resp = &wasmbindingstypes.RequestResponse{
		Code:   paramtypes.RES_OK,
		Result: "Redelegation succeed!",
	}

	return resp, nil
}
