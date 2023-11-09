package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/incentive/keeper"
)

// Messenger handles messages for the Incentive module.
type Messenger struct {
	keeper           *keeper.Keeper
	stakingKeeper    *stakingkeeper.Keeper
	commitmentKeeper *commitmentkeeper.Keeper
}

func NewMessenger(keeper *keeper.Keeper, stakingKeeper *stakingkeeper.Keeper, commitmentKeeper *commitmentkeeper.Keeper) *Messenger {
	return &Messenger{
		keeper:           keeper,
		stakingKeeper:    stakingKeeper,
		commitmentKeeper: commitmentKeeper,
	}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmbindingstypes.ElysMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.MsgBeginRedelegate != nil:
		return m.msgBeginRedelegate(ctx, contractAddr, msg.MsgBeginRedelegate)
	case msg.MsgCancelUnbondingDelegation != nil:
		return m.msgCancelUnbondingDelegation(ctx, contractAddr, msg.MsgCancelUnbondingDelegation)
	case msg.MsgWithdrawRewards != nil:
		return m.msgWithdrawRewards(ctx, contractAddr, msg.MsgWithdrawRewards)
	case msg.MsgWithdrawValidatorCommission != nil:
		return m.msgWithdrawValidatorCommission(ctx, contractAddr, msg.MsgWithdrawValidatorCommission)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
