package abci

import (
	"encoding/json"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/keeper"
	"github.com/elys-network/elys/v7/x/clob/types"
)

type VoteExtensionHandler struct {
	logger log.Logger
	keeper keeper.Keeper
}

func NewVoteExtensionHandler(logger log.Logger, keeper keeper.Keeper) *VoteExtensionHandler {
	return &VoteExtensionHandler{
		logger: logger,
		keeper: keeper,
	}
}

type VoteExtension struct {
	MatchedOrders []types.MatchedOrderExecution `json:"matched_orders"`
	Height        int64                         `json:"height"`
}

func (h *VoteExtensionHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *cometabci.RequestExtendVote) (*cometabci.ResponseExtendVote, error) {
		h.logger.Info("ExtendVote", "height", req.Height)

		// Get all active market IDs
		marketIds := h.getActiveMarketIds(ctx)

		// Get matched orders to propose (limit to 100 matches per market to avoid too large extensions)
		matchedOrders := h.keeper.GetOperationsToPropose(ctx, marketIds, 100)

		// Create vote extension
		voteExt := VoteExtension{
			MatchedOrders: matchedOrders,
			Height:        req.Height,
		}

		// Marshal to JSON
		bz, err := json.Marshal(voteExt)
		if err != nil {
			h.logger.Error("failed to marshal vote extension", "error", err)
			return &cometabci.ResponseExtendVote{}, nil
		}

		h.logger.Info("Vote extension created",
			"height", req.Height,
			"matched_orders", len(matchedOrders))

		return &cometabci.ResponseExtendVote{
			VoteExtension: bz,
		}, nil
	}
}

func (h *VoteExtensionHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *cometabci.RequestVerifyVoteExtension) (*cometabci.ResponseVerifyVoteExtension, error) {
		h.logger.Debug("VerifyVoteExtension", "height", req.Height)

		// Basic validation - just check if it can be unmarshaled
		var voteExt VoteExtension
		if err := json.Unmarshal(req.VoteExtension, &voteExt); err != nil {
			h.logger.Error("failed to unmarshal vote extension", "error", err)
			return &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			}, nil
		}

		// Verify height matches
		if voteExt.Height != req.Height {
			h.logger.Error("vote extension height mismatch",
				"expected", req.Height,
				"got", voteExt.Height)
			return &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			}, nil
		}

		// Basic validation of matched orders
		for _, order := range voteExt.MatchedOrders {
			if order.Quantity.IsNegative() || order.Price.IsNegative() {
				h.logger.Error("invalid matched order", "order", order)
				return &cometabci.ResponseVerifyVoteExtension{
					Status: cometabci.ResponseVerifyVoteExtension_REJECT,
				}, nil
			}
		}

		return &cometabci.ResponseVerifyVoteExtension{
			Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
		}, nil
	}
}

func (h *VoteExtensionHandler) getActiveMarketIds(ctx sdk.Context) []uint64 {
	// Get all perpetual markets
	markets := h.keeper.GetAllPerpetualMarket(ctx)
	marketIds := make([]uint64, 0, len(markets))

	for _, market := range markets {
		// Add all markets since we don't have a Status field
		// You can add filtering logic here if needed
		marketIds = append(marketIds, market.Id)
	}

	return marketIds
}
