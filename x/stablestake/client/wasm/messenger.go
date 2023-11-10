package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/stablestake/keeper"
)

// Messenger handles messages for the Stable Stake module.
type Messenger struct {
	keeper *keeper.Keeper
}

func NewMessenger(keeper *keeper.Keeper) *Messenger {
	return &Messenger{keeper: keeper}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmbindingstypes.ElysMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.StakestakeBond != nil:
		return m.msgBond(ctx, contractAddr, msg.StakestakeBond)
	case msg.StakestakeUnbond != nil:
		return m.msgUnbond(ctx, contractAddr, msg.StakestakeUnbond)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
