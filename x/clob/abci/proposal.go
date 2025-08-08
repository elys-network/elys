package abci

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/elys-network/elys/v7/x/clob/keeper"
	"github.com/elys-network/elys/v7/x/clob/types"
	protov2 "google.golang.org/protobuf/proto"
)

type ProposalHandler struct {
	logger    log.Logger
	keeper    keeper.Keeper
	txEncoder sdk.TxEncoder
	txDecoder sdk.TxDecoder
}

func NewProposalHandler(
	logger log.Logger,
	keeper keeper.Keeper,
	txEncoder sdk.TxEncoder,
	txDecoder sdk.TxDecoder,
) *ProposalHandler {
	return &ProposalHandler{
		logger:    logger,
		keeper:    keeper,
		txEncoder: txEncoder,
		txDecoder: txDecoder,
	}
}

func (h *ProposalHandler) PrepareProposalHandler() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *cometabci.RequestPrepareProposal) (*cometabci.ResponsePrepareProposal, error) {
		h.logger.Info("PrepareProposal", "height", req.Height)

		// Aggregate matched orders from vote extensions
		aggregatedOrders := h.aggregateVoteExtensions(req.LocalLastCommit)

		// Create the ExecuteMatchedOrders message if we have orders to execute
		if len(aggregatedOrders) > 0 {
			// Note: No sender field needed - this is a system-generated message
			msg := &types.MsgExecuteMatchedOrders{
				MatchedOrders: aggregatedOrders,
			}

			// Create a transaction with the message
			tx, err := h.createTx(msg)
			if err != nil {
				h.logger.Error("failed to create tx", "error", err)
				return &cometabci.ResponsePrepareProposal{Txs: req.Txs}, nil
			}

			// Encode the transaction
			txBytes, err := h.txEncoder(tx)
			if err != nil {
				h.logger.Error("failed to encode tx", "error", err)
				return &cometabci.ResponsePrepareProposal{Txs: req.Txs}, nil
			}

			// Prepend our tx to the list
			txs := [][]byte{txBytes}
			txs = append(txs, req.Txs...)

			h.logger.Info("Added matched orders tx to proposal",
				"orders", len(aggregatedOrders),
				"tx_size", len(txBytes))

			return &cometabci.ResponsePrepareProposal{Txs: txs}, nil
		}

		return &cometabci.ResponsePrepareProposal{Txs: req.Txs}, nil
	}
}

func (h *ProposalHandler) ProcessProposalHandler() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *cometabci.RequestProcessProposal) (*cometabci.ResponseProcessProposal, error) {
		h.logger.Debug("ProcessProposal", "height", req.Height, "txs", len(req.Txs))

		// Check if the first tx is our matched orders tx
		if len(req.Txs) > 0 {
			tx, err := h.txDecoder(req.Txs[0])
			if err == nil {
				msgs := tx.GetMsgs()
				if len(msgs) == 1 {
					if matchedOrdersMsg, ok := msgs[0].(*types.MsgExecuteMatchedOrders); ok {
						h.logger.Info("Found matched orders tx in proposal",
							"orders_count", len(matchedOrdersMsg.MatchedOrders))

						// CRITICAL VALIDATION: Verify the proposed orders are valid
						// Since ProcessProposal doesn't have access to vote extensions,
						// we validate against the keeper's current state of matched orders
						// Each validator will check independently based on their local orderbook
						isValid := h.validateMatchedOrdersFromKeeper(
							ctx,
							matchedOrdersMsg.MatchedOrders,
						)

						if !isValid {
							h.logger.Error("Invalid matched orders in proposal - rejecting",
								"height", req.Height,
								"proposed_orders", len(matchedOrdersMsg.MatchedOrders))
							return &cometabci.ResponseProcessProposal{
								Status: cometabci.ResponseProcessProposal_REJECT,
							}, nil
						}

						h.logger.Info("Matched orders validated successfully")
					}
				}
			}
		}

		return &cometabci.ResponseProcessProposal{Status: cometabci.ResponseProcessProposal_ACCEPT}, nil
	}
}

// getOrderKey creates a unique key for an order
func (h *ProposalHandler) getOrderKey(order types.MatchedOrderExecution) string {
	return fmt.Sprintf("%d-%d-%d-%s-%s-%s-%s",
		order.MarketId,
		order.BuyOrderCounter,
		order.SellOrderCounter,
		order.Price.String(),
		order.Quantity.String(),
		order.Buyer,
		order.Seller)
}

// validateMatchedOrdersFromKeeper validates the proposed orders against the keeper's state
// Each validator independently verifies based on their local orderbook view
func (h *ProposalHandler) validateMatchedOrdersFromKeeper(
	ctx sdk.Context,
	proposedOrders []types.MatchedOrderExecution,
) bool {
	// Get the active market IDs
	marketIds := h.getActiveMarketIds(ctx)

	// Get the matched orders from the keeper (what this validator would propose)
	expectedOrders := h.keeper.GetOperationsToPropose(ctx, marketIds, 100)

	// Now validate that proposed orders match expected orders
	return h.validateOrdersMatch(proposedOrders, expectedOrders)
}

// getActiveMarketIds returns all active market IDs
func (h *ProposalHandler) getActiveMarketIds(ctx sdk.Context) []uint64 {
	markets := h.keeper.GetAllPerpetualMarket(ctx)
	marketIds := make([]uint64, 0, len(markets))

	for _, market := range markets {
		marketIds = append(marketIds, market.Id)
	}

	return marketIds
}

// validateOrdersMatch ensures the proposed orders exactly match the expected orders
// This prevents a malicious proposer from injecting fake orders or censoring legitimate ones
func (h *ProposalHandler) validateOrdersMatch(
	proposedOrders []types.MatchedOrderExecution,
	expectedOrders []types.MatchedOrderExecution,
) bool {
	// First check counts match
	if len(proposedOrders) != len(expectedOrders) {
		h.logger.Error("Order count mismatch",
			"proposed_count", len(proposedOrders),
			"expected_count", len(expectedOrders))
		return false
	}

	// Build a map of expected orders for efficient lookup
	expectedMap := make(map[string]types.MatchedOrderExecution)
	for _, order := range expectedOrders {
		key := h.getOrderKey(order)
		expectedMap[key] = order
	}

	// Check that all proposed orders are in the expected set and match exactly
	for _, proposedOrder := range proposedOrders {
		key := h.getOrderKey(proposedOrder)
		expectedOrder, exists := expectedMap[key]
		if !exists {
			h.logger.Error("Proposed order not in expected orders from vote extensions",
				"order_key", key)
			return false
		}

		// Double-check that the order details match exactly
		if !ordersMatch(proposedOrder, expectedOrder) {
			h.logger.Error("Proposed order details don't match expected order",
				"order_key", key)
			return false
		}

		// Remove from map to detect duplicates
		delete(expectedMap, key)
	}

	// Check if any expected orders were not included (should be empty now)
	if len(expectedMap) > 0 {
		h.logger.Error("Some expected orders were not included in proposal",
			"missing_count", len(expectedMap))
		for key := range expectedMap {
			h.logger.Error("Missing order", "order_key", key)
		}
		return false
	}

	h.logger.Info("Order validation successful",
		"orders_count", len(proposedOrders))

	return true
}

// ordersMatch checks if two orders are exactly the same
func ordersMatch(a, b types.MatchedOrderExecution) bool {
	return a.MarketId == b.MarketId &&
		a.BuyOrderCounter == b.BuyOrderCounter &&
		a.SellOrderCounter == b.SellOrderCounter &&
		a.Price.Equal(b.Price) &&
		a.Quantity.Equal(b.Quantity) &&
		a.Buyer == b.Buyer &&
		a.Seller == b.Seller &&
		a.BuyerSubAccountId == b.BuyerSubAccountId &&
		a.SellerSubAccountId == b.SellerSubAccountId
}

func (h *ProposalHandler) aggregateVoteExtensions(extCommit cometabci.ExtendedCommitInfo) []types.MatchedOrderExecution {
	// Use voting power weighted aggregation
	orderVotingPower := make(map[string]int64)
	orderMap := make(map[string]types.MatchedOrderExecution)
	totalVotingPower := int64(0)

	for _, vote := range extCommit.Votes {
		if len(vote.VoteExtension) == 0 || vote.Validator.Power == 0 {
			continue
		}

		totalVotingPower += vote.Validator.Power

		var voteExt VoteExtension
		if err := json.Unmarshal(vote.VoteExtension, &voteExt); err != nil {
			h.logger.Error("failed to unmarshal vote extension", "error", err)
			continue
		}

		// Add matched orders to map with voting power tracking
		for _, order := range voteExt.MatchedOrders {
			key := fmt.Sprintf("%d-%d-%d-%s-%s",
				order.MarketId,
				order.BuyOrderCounter,
				order.SellOrderCounter,
				order.Price.String(),
				order.Quantity.String())
			orderMap[key] = order
			orderVotingPower[key] += vote.Validator.Power
		}
	}

	// Only include orders that have at least 2/3 voting power support
	requiredPower := (totalVotingPower * 2) / 3
	aggregated := make([]types.MatchedOrderExecution, 0)

	for key, order := range orderMap {
		if orderVotingPower[key] >= requiredPower {
			aggregated = append(aggregated, order)
		}
	}

	h.logger.Info("Aggregated vote extensions",
		"votes", len(extCommit.Votes),
		"total_voting_power", totalVotingPower,
		"required_power", requiredPower,
		"matched_orders", len(aggregated))

	return aggregated
}

func (h *ProposalHandler) createTx(msg sdk.Msg) (sdk.Tx, error) {
	// This is a simplified version - in production you'd need proper tx building
	// with gas, fees, signatures, etc.
	return &SimpleTx{msgs: []sdk.Msg{msg}}, nil
}

// SimpleTx is a minimal Tx implementation for vote extension transactions
// This transaction type has no signers and requires no signatures
type SimpleTx struct {
	msgs []sdk.Msg
}

func (tx *SimpleTx) GetMsgs() []sdk.Msg {
	return tx.msgs
}

func (tx *SimpleTx) GetMsgsV2() ([]protov2.Message, error) {
	// Convert sdk.Msg to protov2.Message
	msgs := make([]protov2.Message, len(tx.msgs))
	for i, msg := range tx.msgs {
		// Type assertion to proto message
		if protoMsg, ok := msg.(protov2.Message); ok {
			msgs[i] = protoMsg
		}
	}
	return msgs, nil
}

func (tx *SimpleTx) ValidateBasic() error {
	// Basic validation can be done here
	// Most msgs don't have ValidateBasic anymore in SDK v0.50+
	return nil
}

// GetSigners returns an empty slice as vote extension transactions have no signers
func (tx *SimpleTx) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

// GetSignaturesV2 returns empty as there are no signatures
func (tx *SimpleTx) GetSignaturesV2() ([]signing.SignatureV2, error) {
	return []signing.SignatureV2{}, nil
}
