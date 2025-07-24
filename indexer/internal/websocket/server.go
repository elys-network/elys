package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/elys-network/elys/indexer/internal/cache"
	"github.com/elys-network/elys/indexer/internal/config"
	"github.com/elys-network/elys/indexer/internal/database"
	"github.com/elys-network/elys/indexer/internal/models"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Server struct {
	config   *config.WebSocketConfig
	upgrader websocket.Upgrader
	hub      *Hub
	cache    cache.CacheInterface
	db       *database.Repository
	logger   *zap.Logger
	server   *http.Server
}

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *BroadcastMessage
	mu         sync.RWMutex
}

type Client struct {
	id            string
	conn          *websocket.Conn
	send          chan []byte
	subscriptions map[string]*models.WebSocketSubscription
	hub           *Hub
	logger        *zap.Logger
	mu            sync.RWMutex
}

type BroadcastMessage struct {
	Channel string
	Data    []byte
}

func NewServer(cfg *config.WebSocketConfig, cacheImpl cache.CacheInterface, db *database.Repository, logger *zap.Logger) *Server {
	hub := &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  cfg.ReadBufferSize,
		WriteBufferSize: cfg.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			// TODO: Implement proper origin checking in production
			return true
		},
	}

	return &Server{
		config:   cfg,
		upgrader: upgrader,
		hub:      hub,
		cache:    cacheImpl,
		db:       db,
		logger:   logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	// Start the hub
	go s.hub.run(ctx)

	// Start Redis subscription handler
	go s.handleRedisSubscriptions(ctx)

	// Set up HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWebSocket)
	mux.HandleFunc("/health", s.handleHealth)

	s.server = &http.Server{
		Addr:    s.config.ListenAddr,
		Handler: mux,
	}

	s.logger.Info("WebSocket server starting", zap.String("addr", s.config.ListenAddr))

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("WebSocket server error", zap.Error(err))
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Check if we've reached max connections
	s.hub.mu.RLock()
	if len(s.hub.clients) >= s.config.MaxConnections {
		s.hub.mu.RUnlock()
		http.Error(w, "Maximum connections reached", http.StatusServiceUnavailable)
		return
	}
	s.hub.mu.RUnlock()

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("Failed to upgrade connection", zap.Error(err))
		return
	}

	clientID := r.Header.Get("X-Client-ID")
	if clientID == "" {
		clientID = fmt.Sprintf("client_%d", time.Now().UnixNano())
	}

	client := &Client{
		id:            clientID,
		conn:          conn,
		send:          make(chan []byte, 256),
		subscriptions: make(map[string]*models.WebSocketSubscription),
		hub:           s.hub,
		logger:        s.logger,
	}

	s.hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump(s.config)
	go client.readPump(s.config, s.db, s.cache)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"clients": len(s.hub.clients),
	})
}

func (h *Hub) run(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.id] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				if client.isSubscribedToChannel(message.Channel) {
					select {
					case client.send <- message.Data:
					default:
						// Client's send channel is full, close it
						close(client.send)
						delete(h.clients, client.id)
					}
				}
			}
			h.mu.RUnlock()

		case <-ticker.C:
			// Clean up inactive clients
			h.cleanupInactiveClients()
		}
	}
}

func (h *Hub) cleanupInactiveClients() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for id, client := range h.clients {
		if err := client.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
			delete(h.clients, id)
			close(client.send)
		}
	}
}

func (c *Client) readPump(cfg *config.WebSocketConfig, db *database.Repository, cache cache.CacheInterface) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(cfg.MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(cfg.PongTimeout))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(cfg.PongTimeout))
		db.UpdateSubscriptionPing(context.Background(), c.id)
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("WebSocket error", zap.Error(err))
			}
			break
		}

		var msg models.WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			c.sendError("Invalid message format")
			continue
		}

		switch msg.Type {
		case models.WSMessageTypeSubscribe:
			c.handleSubscribe(msg, db, cache)
		case models.WSMessageTypeUnsubscribe:
			c.handleUnsubscribe(msg, db)
		case models.WSMessageTypePing:
			c.handlePing()
		}
	}
}

func (c *Client) writePump(cfg *config.WebSocketConfig) {
	ticker := time.NewTicker(cfg.PingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(cfg.WriteTimeout))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(cfg.WriteTimeout))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleSubscribe(msg models.WSMessage, db *database.Repository, cache cache.CacheInterface) {
	var req models.WSSubscribeRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		c.sendError("Invalid subscribe request")
		return
	}

	// Create subscription
	sub := &models.WebSocketSubscription{
		ID:               fmt.Sprintf("%s_%s_%d", c.id, req.Type, time.Now().UnixNano()),
		ClientID:         c.id,
		SubscriptionType: string(req.Type),
		Filters:          req.Filters,
		CreatedAt:        time.Now(),
		LastPing:         time.Now(),
	}

	// Store in database
	if err := db.CreateSubscription(context.Background(), sub); err != nil {
		c.logger.Error("Failed to create subscription", zap.Error(err))
		c.sendError("Failed to create subscription")
		return
	}

	// Store in cache
	if err := cache.AddSubscription(context.Background(), c.id, sub); err != nil {
		c.logger.Error("Failed to cache subscription", zap.Error(err))
	}

	// Store locally
	c.mu.Lock()
	c.subscriptions[sub.ID] = sub
	c.mu.Unlock()

	// Send confirmation
	c.sendMessage(models.WSMessage{
		Type:    models.WSMessageTypeUpdate,
		Channel: "subscription",
		Data:    json.RawMessage(`{"status":"subscribed","id":"` + sub.ID + `"}`),
	})
}

func (c *Client) handleUnsubscribe(msg models.WSMessage, db *database.Repository) {
	var data map[string]string
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		c.sendError("Invalid unsubscribe request")
		return
	}

	subID := data["id"]
	if subID == "" {
		c.sendError("Missing subscription ID")
		return
	}

	// Remove from database
	if err := db.DeleteSubscription(context.Background(), subID); err != nil {
		c.logger.Error("Failed to delete subscription", zap.Error(err))
	}

	// Remove locally
	c.mu.Lock()
	delete(c.subscriptions, subID)
	c.mu.Unlock()

	// Send confirmation
	c.sendMessage(models.WSMessage{
		Type:    models.WSMessageTypeUpdate,
		Channel: "subscription",
		Data:    json.RawMessage(`{"status":"unsubscribed","id":"` + subID + `"}`),
	})
}

func (c *Client) handlePing() {
	c.sendMessage(models.WSMessage{
		Type: models.WSMessageTypePong,
	})
}

func (c *Client) sendMessage(msg models.WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		c.logger.Error("Failed to marshal message", zap.Error(err))
		return
	}

	select {
	case c.send <- data:
	default:
		c.logger.Warn("Client send channel full", zap.String("client_id", c.id))
	}
}

func (c *Client) sendError(errMsg string) {
	c.sendMessage(models.WSMessage{
		Type:  models.WSMessageTypeError,
		Error: errMsg,
	})
}

func (c *Client) isSubscribedToChannel(channel string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, sub := range c.subscriptions {
		// Check if subscription matches the channel
		// This is a simplified version - you might want more sophisticated matching
		switch models.WSSubscriptionType(sub.SubscriptionType) {
		case models.WSSubscriptionOrderBook:
			if channel == "orderbook" {
				return true
			}
		case models.WSSubscriptionTrades:
			if channel == "trades" {
				return true
			}
		case models.WSSubscriptionOrders:
			if channel == "orders:"+c.id {
				return true
			}
		case models.WSSubscriptionPositions:
			if channel == "positions:"+c.id {
				return true
			}
		}
	}

	return false
}

func (s *Server) handleRedisSubscriptions(ctx context.Context) {
	// Subscribe to all relevant Redis channels
	pubsub := s.cache.Subscribe(ctx, "orderbook:*", "trades:*", "orders:*", "positions:*")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			// Broadcast to relevant clients
			s.hub.broadcast <- &BroadcastMessage{
				Channel: msg.Channel,
				Data:    []byte(msg.Payload),
			}
		}
	}
}
