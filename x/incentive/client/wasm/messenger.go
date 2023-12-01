package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/incentive/keeper"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
)

// Messenger handles messages for the Incentive module.
type Messenger struct {
	keeper           *keeper.Keeper
	stakingKeeper    *stakingkeeper.Keeper
	commitmentKeeper *commitmentkeeper.Keeper
	parameterKeeper  *parameterkeeper.Keeper
}

func NewMessenger(
	keeper *keeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
	commitmentKeeper *commitmentkeeper.Keeper,
	parameterKeeper *parameterkeeper.Keeper,
) *Messenger {
	return &Messenger{
		keeper:           keeper,
		stakingKeeper:    stakingKeeper,
		commitmentKeeper: commitmentKeeper,
	}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmbindingstypes.ElysMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.IncentiveBeginRedelegate != nil:
		return m.msgBeginRedelegate(ctx, contractAddr, msg.IncentiveBeginRedelegate)
	case msg.IncentiveCancelUnbondingDelegation != nil:
		return m.msgCancelUnbondingDelegation(ctx, contractAddr, msg.IncentiveCancelUnbondingDelegation)
	case msg.IncentiveWithdrawRewards != nil:
		return m.msgWithdrawRewards(ctx, contractAddr, msg.IncentiveWithdrawRewards)
	case msg.IncentiveWithdrawValidatorCommission != nil:
		return m.msgWithdrawValidatorCommission(ctx, contractAddr, msg.IncentiveWithdrawValidatorCommission)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
