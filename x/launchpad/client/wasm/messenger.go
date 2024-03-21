package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/launchpad/keeper"
	parameterkeeper "github.com/elys-network/elys/x/parameter/keeper"
)

// Messenger handles messages for the Leverage LP module.
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
	case msg.LaunchpadBuyElys != nil:
		return m.msgBuyElys(ctx, contractAddr, msg.LaunchpadBuyElys)
	case msg.LaunchpadReturnElys != nil:
		return m.msgReturnElys(ctx, contractAddr, msg.LaunchpadReturnElys)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
