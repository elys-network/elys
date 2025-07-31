package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/elys-network/elys/indexer/internal/cache"
	"github.com/elys-network/elys/indexer/internal/database"
	"github.com/elys-network/elys/indexer/internal/models"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// CLOBHandler handles CLOB-related API endpoints
type CLOBHandler struct {
	db     *database.Repository
	cache  *cache.Cache
	logger *zap.Logger
}

func NewCLOBHandler(db *database.Repository, cache *cache.Cache, logger *zap.Logger) *CLOBHandler {
	return &CLOBHandler{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

// RegisterRoutes registers all CLOB API routes
func (h *CLOBHandler) RegisterRoutes(router *mux.Router) {
	// Market endpoints
	router.HandleFunc("/api/v1/clob/markets", h.GetMarkets).Methods("GET")
	router.HandleFunc("/api/v1/clob/markets/{marketId}", h.GetMarket).Methods("GET")
	router.HandleFunc("/api/v1/clob/markets/{marketId}/stats", h.GetMarketStats).Methods("GET")

	// Order book endpoints
	router.HandleFunc("/api/v1/clob/markets/{marketId}/orderbook", h.GetOrderBook).Methods("GET")
	router.HandleFunc("/api/v1/clob/markets/{marketId}/orderbook/snapshot", h.GetOrderBookSnapshot).Methods("GET")
	router.HandleFunc("/api/v1/clob/markets/{marketId}/best", h.GetBestBidAsk).Methods("GET")

	// Order endpoints
	router.HandleFunc("/api/v1/clob/orders/active", h.GetActiveOrders).Methods("GET")
	router.HandleFunc("/api/v1/clob/orders/{orderId}", h.GetOrder).Methods("GET")
	router.HandleFunc("/api/v1/clob/users/{address}/orders", h.GetUserOrders).Methods("GET")

	// Trade endpoints
	router.HandleFunc("/api/v1/clob/markets/{marketId}/trades", h.GetMarketTrades).Methods("GET")
	router.HandleFunc("/api/v1/clob/users/{address}/trades", h.GetUserTrades).Methods("GET")

	// Position endpoints
	router.HandleFunc("/api/v1/clob/users/{address}/positions", h.GetUserPositions).Methods("GET")
	router.HandleFunc("/api/v1/clob/positions/{positionId}", h.GetPosition).Methods("GET")

	// Funding rate endpoints
	router.HandleFunc("/api/v1/clob/markets/{marketId}/funding", h.GetFundingRates).Methods("GET")

	// Order matching support endpoints
	router.HandleFunc("/api/v1/clob/matching/orders", h.GetOrdersForMatching).Methods("GET")
	router.HandleFunc("/api/v1/clob/matching/lock", h.LockOrderForMatching).Methods("POST")
	router.HandleFunc("/api/v1/clob/matching/unlock", h.UnlockOrder).Methods("POST")
}

// GetMarkets returns all CLOB markets
func (h *CLOBHandler) GetMarkets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	markets, err := h.db.GetActiveCLOBMarkets(ctx)
	if err != nil {
		h.writeError(w, "Failed to get markets", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, markets)
}

// GetMarket returns a specific market
func (h *CLOBHandler) GetMarket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	marketID, err := strconv.ParseUint(vars["marketId"], 10, 64)
	if err != nil {
		h.writeError(w, "Invalid market ID", http.StatusBadRequest)
		return
	}

	market, err := h.db.GetCLOBMarket(ctx, marketID)
	if err != nil {
		h.writeError(w, "Market not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, market)
}

// GetMarketStats returns market statistics
func (h *CLOBHandler) GetMarketStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	marketID, err := strconv.ParseUint(vars["marketId"], 10, 64)
	if err != nil {
		h.writeError(w, "Invalid market ID", http.StatusBadRequest)
		return
	}

	stats, err := h.db.GetCLOBMarketStats(ctx, marketID)
	if err != nil {
		h.writeError(w, "Failed to get market stats", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, stats)
}

// GetOrderBook returns the current order book from cache
func (h *CLOBHandler) GetOrderBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	marketID, err := strconv.ParseUint(vars["marketId"], 10, 64)
	if err != nil {
		h.writeError(w, "Invalid market ID", http.StatusBadRequest)
		return
	}

	// Get depth limit from query params
	depthStr := r.URL.Query().Get("depth")
	depth := 20 // default depth
	if depthStr != "" {
		if d, err := strconv.Atoi(depthStr); err == nil && d > 0 && d <= 100 {
			depth = d
		}
	}

	// Get active orders
	orders, err := h.db.GetActiveOrdersForMarket(ctx, marketID)
	if err != nil {
		h.writeError(w, "Failed to get order book", http.StatusInternalServerError)
		return
	}

	// Build order book
	orderBook := h.buildOrderBook(orders, depth)
	h.writeJSON(w, orderBook)
}

// GetOrderBookSnapshot returns the latest order book snapshot
func (h *CLOBHandler) GetOrderBookSnapshot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	marketID, err := strconv.ParseUint(vars["marketId"], 10, 64)
	if err != nil {
		h.writeError(w, "Invalid market ID", http.StatusBadRequest)
		return
	}

	snapshot, err := h.db.GetLatestCLOBOrderBookSnapshot(ctx, marketID)
	if err != nil {
		h.writeError(w, "Failed to get order book snapshot", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, snapshot)
}

// GetBestBidAsk returns the best bid and ask prices
func (h *CLOBHandler) GetBestBidAsk(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	marketID, err := strconv.ParseUint(vars["marketId"], 10, 64)
	if err != nil {
		h.writeError(w, "Invalid market ID", http.StatusBadRequest)
		return
	}

	bestBid, bestAsk, err := h.cache.GetBestBidAsk(ctx, marketID)
	if err != nil {
		h.writeError(w, "Failed to get best bid/ask", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"market_id": marketID,
		"best_bid":  nil,
		"best_ask":  nil,
		"spread":    nil,
		"mid_price": nil,
		"timestamp": time.Now(),
	}

	if bestBid != nil {
		response["best_bid"] = bestBid.String()
	}
	if bestAsk != nil {
		response["best_ask"] = bestAsk.String()
	}
	if bestBid != nil && bestAsk != nil {
		spread := bestAsk.Sub(*bestBid)
		midPrice := bestBid.Add(*bestAsk).Div(decimal.NewFromInt(2))
		response["spread"] = spread.String()
		response["mid_price"] = midPrice.String()
	}

	h.writeJSON(w, response)
}

// GetActiveOrders returns all active orders across markets
func (h *CLOBHandler) GetActiveOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get query parameters
	marketIDStr := r.URL.Query().Get("market_id")
	owner := r.URL.Query().Get("owner")
	orderType := r.URL.Query().Get("order_type")
	limitStr := r.URL.Query().Get("limit")

	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	// Build query based on filters
	var orders []*models.CLOBOrder
	var err error

	if owner != "" {
		statuses := []models.CLOBOrderStatus{
			models.CLOBOrderStatusPending,
			models.CLOBOrderStatusPartiallyFilled,
		}
		orders, err = h.db.GetCLOBOrdersByOwner(ctx, owner, statuses, limit)
	} else if marketIDStr != "" {
		marketID, err := strconv.ParseUint(marketIDStr, 10, 64)
		if err != nil {
			h.writeError(w, "Invalid market ID", http.StatusBadRequest)
			return
		}
		orders, err = h.db.GetActiveOrdersForMarket(ctx, marketID)
	} else {
		h.writeError(w, "Must specify either market_id or owner", http.StatusBadRequest)
		return
	}

	if err != nil {
		h.writeError(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	// Filter by order type if specified
	if orderType != "" {
		filtered := make([]*models.CLOBOrder, 0)
		for _, order := range orders {
			if string(order.OrderType) == orderType {
				filtered = append(filtered, order)
			}
		}
		orders = filtered
	}

	h.writeJSON(w, orders)
}

// GetOrder returns a specific order
func (h *CLOBHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	orderID, err := strconv.ParseUint(vars["orderId"], 10, 64)
	if err != nil {
		h.writeError(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.db.GetCLOBOrder(ctx, orderID)
	if err != nil {
		h.writeError(w, "Order not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, order)
}

// GetUserOrders returns orders for a specific user
func (h *CLOBHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	address := vars["address"]
	statusStr := r.URL.Query().Get("status")
	limitStr := r.URL.Query().Get("limit")

	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	var statuses []models.CLOBOrderStatus
	if statusStr != "" {
		statuses = append(statuses, models.CLOBOrderStatus(statusStr))
	}

	orders, err := h.db.GetCLOBOrdersByOwner(ctx, address, statuses, limit)
	if err != nil {
		h.writeError(w, "Failed to get user orders", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, orders)
}

// GetMarketTrades returns recent trades for a market
func (h *CLOBHandler) GetMarketTrades(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	marketID, err := strconv.ParseUint(vars["marketId"], 10, 64)
	if err != nil {
		h.writeError(w, "Invalid market ID", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	trades, err := h.db.GetCLOBTradeHistory(ctx, marketID, "", limit)
	if err != nil {
		h.writeError(w, "Failed to get trades", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, trades)
}

// GetUserTrades returns trades for a specific user
func (h *CLOBHandler) GetUserTrades(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	address := vars["address"]
	marketIDStr := r.URL.Query().Get("market_id")
	limitStr := r.URL.Query().Get("limit")

	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	var marketID uint64
	if marketIDStr != "" {
		if id, err := strconv.ParseUint(marketIDStr, 10, 64); err == nil {
			marketID = id
		}
	}

	trades, err := h.db.GetCLOBTradeHistory(ctx, marketID, address, limit)
	if err != nil {
		h.writeError(w, "Failed to get user trades", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, trades)
}

// GetUserPositions returns open positions for a user
func (h *CLOBHandler) GetUserPositions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	address := vars["address"]
	marketIDStr := r.URL.Query().Get("market_id")

	var marketID uint64
	if marketIDStr != "" {
		if id, err := strconv.ParseUint(marketIDStr, 10, 64); err == nil {
			marketID = id
		}
	}

	positions, err := h.db.GetOpenCLOBPositions(ctx, address, marketID)
	if err != nil {
		h.writeError(w, "Failed to get user positions", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, positions)
}

// GetPosition returns a specific position
func (h *CLOBHandler) GetPosition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	positionID, err := strconv.ParseUint(vars["positionId"], 10, 64)
	if err != nil {
		h.writeError(w, "Invalid position ID", http.StatusBadRequest)
		return
	}

	position, err := h.db.GetCLOBPosition(ctx, positionID)
	if err != nil {
		h.writeError(w, "Position not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, position)
}

// GetFundingRates returns funding rate history
func (h *CLOBHandler) GetFundingRates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	marketID, err := strconv.ParseUint(vars["marketId"], 10, 64)
	if err != nil {
		h.writeError(w, "Invalid market ID", http.StatusBadRequest)
		return
	}

	// TODO: Implement GetCLOBFundingRates in repository
	h.writeJSON(w, []interface{}{})
}

// Order Matching Support Endpoints

// GetOrdersForMatching returns orders suitable for matching
func (h *CLOBHandler) GetOrdersForMatching(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	marketIDStr := r.URL.Query().Get("market_id")
	if marketIDStr == "" {
		h.writeError(w, "market_id is required", http.StatusBadRequest)
		return
	}

	marketID, err := strconv.ParseUint(marketIDStr, 10, 64)
	if err != nil {
		h.writeError(w, "Invalid market ID", http.StatusBadRequest)
		return
	}

	// Get active orders
	orders, err := h.db.GetActiveOrdersForMarket(ctx, marketID)
	if err != nil {
		h.writeError(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	// Group by order type
	response := map[string]interface{}{
		"market_id":   marketID,
		"timestamp":   time.Now(),
		"buy_orders":  make([]*models.CLOBOrder, 0),
		"sell_orders": make([]*models.CLOBOrder, 0),
	}

	for _, order := range orders {
		if order.OrderType == models.CLOBOrderTypeLimitBuy {
			response["buy_orders"] = append(response["buy_orders"].([]*models.CLOBOrder), order)
		} else if order.OrderType == models.CLOBOrderTypeLimitSell {
			response["sell_orders"] = append(response["sell_orders"].([]*models.CLOBOrder), order)
		}
	}

	h.writeJSON(w, response)
}

// LockOrderForMatching locks an order for matching
func (h *CLOBHandler) LockOrderForMatching(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		OrderID          uint64 `json:"order_id"`
		MatchingEngineID string `json:"matching_engine_id"`
		TTL              int    `json:"ttl_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.OrderID == 0 || req.MatchingEngineID == "" {
		h.writeError(w, "order_id and matching_engine_id are required", http.StatusBadRequest)
		return
	}

	ttl := time.Duration(req.TTL) * time.Second
	if ttl == 0 {
		ttl = 5 * time.Second // default 5 seconds
	}

	locked, err := h.cache.LockCLOBOrder(ctx, req.OrderID, req.MatchingEngineID, ttl)
	if err != nil {
		h.writeError(w, "Failed to lock order", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"locked":   locked,
		"order_id": req.OrderID,
		"ttl":      ttl.Seconds(),
	}

	h.writeJSON(w, response)
}

// UnlockOrder unlocks an order
func (h *CLOBHandler) UnlockOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		OrderID          uint64 `json:"order_id"`
		MatchingEngineID string `json:"matching_engine_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.OrderID == 0 || req.MatchingEngineID == "" {
		h.writeError(w, "order_id and matching_engine_id are required", http.StatusBadRequest)
		return
	}

	err := h.cache.UnlockCLOBOrder(ctx, req.OrderID, req.MatchingEngineID)
	if err != nil {
		h.writeError(w, "Failed to unlock order", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"unlocked": true,
		"order_id": req.OrderID,
	}

	h.writeJSON(w, response)
}

// Helper methods

func (h *CLOBHandler) buildOrderBook(orders []*models.CLOBOrder, depth int) map[string]interface{} {
	bids := make([]map[string]string, 0)
	asks := make([]map[string]string, 0)

	bidsByPrice := make(map[string]decimal.Decimal)
	asksByPrice := make(map[string]decimal.Decimal)

	for _, order := range orders {
		priceStr := order.Price.String()

		if order.OrderType == models.CLOBOrderTypeLimitBuy {
			bidsByPrice[priceStr] = bidsByPrice[priceStr].Add(order.RemainingAmount)
		} else if order.OrderType == models.CLOBOrderTypeLimitSell {
			asksByPrice[priceStr] = asksByPrice[priceStr].Add(order.RemainingAmount)
		}
	}

	// Convert to sorted arrays
	for price, amount := range bidsByPrice {
		bids = append(bids, map[string]string{
			"price":  price,
			"amount": amount.String(),
		})
	}

	for price, amount := range asksByPrice {
		asks = append(asks, map[string]string{
			"price":  price,
			"amount": amount.String(),
		})
	}

	// Sort and limit depth
	// TODO: Implement proper sorting

	return map[string]interface{}{
		"bids":      bids,
		"asks":      asks,
		"timestamp": time.Now(),
	}
}

func (h *CLOBHandler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *CLOBHandler) writeError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
