package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/tradeshield/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {

	// Testnet
	if ctx.ChainID() == "elysicstestnet-1" {
		allOrders := m.keeper.GetAllPendingSpotOrder(ctx)
		for _, order := range allOrders {
			if order.OrderType == types.SpotOrderType_LIMITBUY {
				m.keeper.RemovePendingSpotOrder(ctx, order.OrderId)
			}
		}
	} else { // Mainnet
		// This is in reference to https://elys.atlassian.net/jira/software/projects/DEV/boards/2?assignee=712020%3A1d0833ef-4e27-4457-98ae-bd3075185694&selectedIssue=DEV-2322
		// Buy orders were not fulfilled due to a issue in arguments passing, although that should have thrown error but
		// error handling was not done so funds were stuck in the pool. This task was completed in PR #1081, so these orders needs to be
		// removed
		ordersIds := []uint64{164, 177, 180}
		for _, orderId := range ordersIds {
			m.keeper.RemovePendingSpotOrder(ctx, orderId)
		}
	}

	return nil
}
