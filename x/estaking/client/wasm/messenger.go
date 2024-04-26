package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/estaking/keeper"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
)

// Messenger handles messages for the Masterchef module.
type Messenger struct {
	keeper          *keeper.Keeper
	parameterKeeper *parameterkeeper.Keeper
}

func NewMessenger(
	keeper *keeper.Keeper,
	parameterKeeper *parameterkeeper.Keeper,
) *Messenger {
	return &Messenger{
		keeper:          keeper,
		parameterKeeper: parameterKeeper,
	}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmbindingstypes.ElysMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.EstakingWithdrawReward != nil:
		return m.msgWithdrawReward(ctx, contractAddr, msg.EstakingWithdrawReward)
	case msg.EstakingWithdrawElysStakingRewards != nil:
		return m.msgWithdrawElysStakingRewards(ctx, contractAddr, msg.EstakingWithdrawElysStakingRewards)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
