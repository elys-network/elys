package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/amm/keeper"
)

// Messenger handles messages for the AMM module.
type Messenger struct {
	keeper *keeper.Keeper
}

func NewMessenger(keeper *keeper.Keeper) *Messenger {
	return &Messenger{keeper: keeper}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmbindingstypes.ElysMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.MsgSwapExactAmountIn != nil:
		return m.msgSwapExactAmountIn(ctx, contractAddr, msg.MsgSwapExactAmountIn)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
