package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/commitment/keeper"
)

// Messenger handles messages for the Commitment module.
type Messenger struct {
	keeper        *keeper.Keeper
	stakingKeeper *stakingkeeper.Keeper
}

func NewMessenger(keeper *keeper.Keeper, stakingKeeper *stakingkeeper.Keeper) *Messenger {
	return &Messenger{
		keeper:        keeper,
		stakingKeeper: stakingKeeper,
	}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmbindingstypes.ElysMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.MsgOpen != nil:
		return m.msgStake(ctx, contractAddr, msg.MsgStake)
	case msg.MsgOpen != nil:
		return m.msgUnstake(ctx, contractAddr, msg.MsgUnstake)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
