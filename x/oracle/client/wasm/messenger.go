package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wasmbindingstypes "github.com/elys-network/elys/wasmbindings/types"
	"github.com/elys-network/elys/x/oracle/keeper"
)

// Messenger handles messages for the Oracle module.
type Messenger struct {
	keeper *keeper.Keeper
}

func NewMessenger(keeper *keeper.Keeper) *Messenger {
	return &Messenger{keeper: keeper}
}

func (m *Messenger) HandleMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmbindingstypes.ElysMsg) ([]sdk.Event, [][]byte, error) {
	switch {
	case msg.OracleFeedPrice != nil:
		return m.msgFeedPrice(ctx, contractAddr, msg.OracleFeedPrice)
	case msg.OracleFeedMultiplePrices != nil:
		return m.msgFeedMultiplePrices(ctx, contractAddr, msg.OracleFeedMultiplePrices)
	case msg.OracleRequestBandPrice != nil:
		return m.msgRequestBandPrice(ctx, contractAddr, msg.OracleRequestBandPrice)
	case msg.OracleSetPriceFeeder != nil:
		return m.msgSetPriceFeeder(ctx, contractAddr, msg.OracleSetPriceFeeder)
	case msg.OracleDeletePriceFeeder != nil:
		return m.msgDeletePriceFeeder(ctx, contractAddr, msg.OracleDeletePriceFeeder)
	default:
		// This handler cannot handle the message
		return nil, nil, wasmbindingstypes.ErrCannotHandleMsg
	}
}
