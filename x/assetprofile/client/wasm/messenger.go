package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/assetprofile/keeper"
)

// Messenger handles messages for the Asset Profile module.
type Messenger struct {
	keeper *keeper.Keeper
}

func NewMessenger(keeper *keeper.Keeper) *Messenger {
	return &Messenger{keeper: keeper}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmbindingstypes.ElysMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.AssetProfileCreateEntry != nil:
		return m.msgCreateEntry(ctx, contractAddr, msg.AssetProfileCreateEntry)
	case msg.AssetProfileUpdateEntry != nil:
		return m.msgUpdateEntry(ctx, contractAddr, msg.AssetProfileUpdateEntry)
	case msg.AssetProfileDeleteEntry != nil:
		return m.msgDeleteEntry(ctx, contractAddr, msg.AssetProfileDeleteEntry)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
