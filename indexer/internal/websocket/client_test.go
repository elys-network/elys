package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elys-network/elys/indexer/internal/models"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWebSocketClient provides a test client for WebSocket testing
type TestWebSocketClient struct {
	conn          *websocket.Conn
	url           string
	receivedMsgs  []models.WSMessage
	mu            sync.RWMutex
	done          chan struct{}
	subscriptions map[string]string // subscription type -> subscription ID
}

func NewTestWebSocketClient(serverURL string) *TestWebSocketClient {
	wsURL := strings.Replace(serverURL, "http", "ws", 1) + "/ws"
	return &TestWebSocketClient{
		url:           wsURL,
		receivedMsgs:  make([]models.WSMessage, 0),
		done:          make(chan struct{}),
		subscriptions: make(map[string]string),
	}
}

func (c *TestWebSocketClient) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	
	// Start reading messages
	go c.readMessages()
	
	return nil
}

func (c *TestWebSocketClient) ConnectWithClientID(clientID string) error {
	header := http.Header{}
	header.Add("X-Client-ID", clientID)
	
	conn, _, err := websocket.DefaultDialer.Dial(c.url, header)
	if err != nil {
		return err
	}
	c.conn = conn
	
	// Start reading messages
	go c.readMessages()
	
	return nil
}

func (c *TestWebSocketClient) Close() error {
	close(c.done)
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *TestWebSocketClient) readMessages() {
	for {
		select {
		case <-c.done:
			return
		default:
			var msg models.WSMessage
			err := c.conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("WebSocket read error: %v\n", err)
				}
				return
			}
			
			c.mu.Lock()
			c.receivedMsgs = append(c.receivedMsgs, msg)
			
			// Track subscription confirmations
			if msg.Type == models.WSMessageTypeUpdate && msg.Channel == "subscription" {
				var data map[string]string
				if err := json.Unmarshal(msg.Data, &data); err == nil {
					if data["status"] == "subscribed" {
						// Extract subscription type from the response
						c.subscriptions[data["type"]] = data["id"]
					}
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *TestWebSocketClient) SendMessage(msg models.WSMessage) error {
	return c.conn.WriteJSON(msg)
}

func (c *TestWebSocketClient) Subscribe(subType models.WSSubscriptionType, filters map[string]interface{}) error {
	req := models.WSSubscribeRequest{
		Type:    subType,
		Filters: filters,
	}
	
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	
	msg := models.WSMessage{
		Type: models.WSMessageTypeSubscribe,
		Data: json.RawMessage(data),
	}
	
	return c.SendMessage(msg)
}

func (c *TestWebSocketClient) Unsubscribe(subscriptionID string) error {
	data, err := json.Marshal(map[string]string{"id": subscriptionID})
	if err != nil {
		return err
	}
	
	msg := models.WSMessage{
		Type: models.WSMessageTypeUnsubscribe,
		Data: json.RawMessage(data),
	}
	
	return c.SendMessage(msg)
}

func (c *TestWebSocketClient) Ping() error {
	msg := models.WSMessage{
		Type: models.WSMessageTypePing,
	}
	return c.SendMessage(msg)
}

func (c *TestWebSocketClient) GetReceivedMessages() []models.WSMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	msgs := make([]models.WSMessage, len(c.receivedMsgs))
	copy(msgs, c.receivedMsgs)
	return msgs
}

func (c *TestWebSocketClient) WaitForMessage(msgType models.WSMessageType, timeout time.Duration) (*models.WSMessage, error) {
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		c.mu.RLock()
		for _, msg := range c.receivedMsgs {
			if msg.Type == msgType {
				c.mu.RUnlock()
				return &msg, nil
			}
		}
		c.mu.RUnlock()
		
		time.Sleep(10 * time.Millisecond)
	}
	
	return nil, fmt.Errorf("timeout waiting for message type: %s", msgType)
}

func (c *TestWebSocketClient) WaitForSubscriptionConfirmation(timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		c.mu.RLock()
		for _, msg := range c.receivedMsgs {
			if msg.Type == models.WSMessageTypeUpdate && msg.Channel == "subscription" {
				var data map[string]string
				if err := json.Unmarshal(msg.Data, &data); err == nil {
					if data["status"] == "subscribed" {
						c.mu.RUnlock()
						return data["id"], nil
					}
				}
			}
		}
		c.mu.RUnlock()
		
		time.Sleep(10 * time.Millisecond)
	}
	
	return "", fmt.Errorf("timeout waiting for subscription confirmation")
}

func (c *TestWebSocketClient) ClearMessages() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.receivedMsgs = c.receivedMsgs[:0]
}

// Test cases using the client

func TestWebSocketClientSubscriptions(t *testing.T) {
	_, httpServer, _, mock := setupTestServer(t)
	defer httpServer.Close()
	
	// Create test client
	client := NewTestWebSocketClient(httpServer.URL)
	err := client.Connect()
	require.NoError(t, err)
	defer client.Close()
	
	// Test OrderBook subscription
	mock.ExpectExec("INSERT INTO websocket_subscriptions").
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	err = client.Subscribe(models.WSSubscriptionOrderBook, map[string]interface{}{
		"pair": "ATOM/USDC",
	})
	require.NoError(t, err)
	
	// Wait for confirmation
	subID, err := client.WaitForSubscriptionConfirmation(2 * time.Second)
	require.NoError(t, err)
	assert.NotEmpty(t, subID)
	
	// Test Orders subscription
	mock.ExpectExec("INSERT INTO websocket_subscriptions").
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	err = client.Subscribe(models.WSSubscriptionOrders, map[string]interface{}{
		"owner": "elys1test123",
	})
	require.NoError(t, err)
	
	subID2, err := client.WaitForSubscriptionConfirmation(2 * time.Second)
	require.NoError(t, err)
	assert.NotEmpty(t, subID2)
	assert.NotEqual(t, subID, subID2)
	
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWebSocketClientBroadcastReceive(t *testing.T) {
	server, httpServer, _, mock := setupTestServer(t)
	defer httpServer.Close()
	
	// Create test client
	client := NewTestWebSocketClient(httpServer.URL)
	err := client.ConnectWithClientID("test-client-123")
	require.NoError(t, err)
	defer client.Close()
	
	// Subscribe to orderbook
	mock.ExpectExec("INSERT INTO websocket_subscriptions").
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	err = client.Subscribe(models.WSSubscriptionOrderBook, map[string]interface{}{})
	require.NoError(t, err)
	
	// Wait for subscription
	_, err = client.WaitForSubscriptionConfirmation(2 * time.Second)
	require.NoError(t, err)
	
	// Clear messages to prepare for broadcast test
	client.ClearMessages()
	
	// Simulate broadcast
	broadcastData := map[string]interface{}{
		"type": "snapshot",
		"data": map[string]interface{}{
			"bids": []map[string]string{
				{"price": "10.00", "amount": "100"},
			},
			"asks": []map[string]string{
				{"price": "10.10", "amount": "50"},
			},
		},
	}
	
	jsonData, _ := json.Marshal(broadcastData)
	server.hub.broadcast <- &BroadcastMessage{
		Channel: "orderbook",
		Data:    jsonData,
	}
	
	// Wait for broadcast message
	time.Sleep(100 * time.Millisecond)
	
	messages := client.GetReceivedMessages()
	assert.Greater(t, len(messages), 0)
	
	// Verify broadcast was received
	var found bool
	for _, msg := range messages {
		var data map[string]interface{}
		if err := json.Unmarshal(msg.Data, &data); err == nil {
			if data["type"] == "snapshot" {
				found = true
				break
			}
		}
	}
	assert.True(t, found, "Broadcast message not received")
}

func TestWebSocketClientReconnection(t *testing.T) {
	_, httpServer, _, _ := setupTestServer(t)
	defer httpServer.Close()
	
	// Create test client
	client := NewTestWebSocketClient(httpServer.URL)
	
	// Connect and disconnect multiple times
	for i := 0; i < 3; i++ {
		err := client.Connect()
		require.NoError(t, err)
		
		// Send ping to verify connection
		err = client.Ping()
		require.NoError(t, err)
		
		// Wait for pong
		pong, err := client.WaitForMessage(models.WSMessageTypePong, 1*time.Second)
		require.NoError(t, err)
		assert.Equal(t, models.WSMessageTypePong, pong.Type)
		
		// Close connection
		err = client.Close()
		require.NoError(t, err)
		
		// Reset for next iteration
		client = NewTestWebSocketClient(httpServer.URL)
	}
}

func TestWebSocketClientErrorHandling(t *testing.T) {
	_, httpServer, _, _ := setupTestServer(t)
	defer httpServer.Close()
	
	// Create test client
	client := NewTestWebSocketClient(httpServer.URL)
	err := client.Connect()
	require.NoError(t, err)
	defer client.Close()
	
	// Send invalid message format
	invalidMsg := []byte("invalid json")
	err = client.conn.WriteMessage(websocket.TextMessage, invalidMsg)
	require.NoError(t, err)
	
	// Wait for error response
	errorMsg, err := client.WaitForMessage(models.WSMessageTypeError, 1*time.Second)
	require.NoError(t, err)
	assert.Equal(t, "Invalid message format", errorMsg.Error)
}

func TestMultipleClientsSimultaneous(t *testing.T) {
	server, httpServer, _, mock := setupTestServer(t)
	defer httpServer.Close()
	
	numClients := 5
	clients := make([]*TestWebSocketClient, numClients)
	
	// Set up mock expectations for all subscriptions
	for i := 0; i < numClients; i++ {
		mock.ExpectExec("INSERT INTO websocket_subscriptions").
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
	
	// Create and connect multiple clients
	var wg sync.WaitGroup
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			
			client := NewTestWebSocketClient(httpServer.URL)
			err := client.ConnectWithClientID(fmt.Sprintf("client-%d", idx))
			require.NoError(t, err)
			
			clients[idx] = client
			
			// Each client subscribes to orderbook
			err = client.Subscribe(models.WSSubscriptionOrderBook, map[string]interface{}{})
			require.NoError(t, err)
			
			// Wait for confirmation
			_, err = client.WaitForSubscriptionConfirmation(2 * time.Second)
			require.NoError(t, err)
		}(i)
	}
	
	wg.Wait()
	
	// Clear messages
	for _, client := range clients {
		client.ClearMessages()
	}
	
	// Broadcast to all clients
	broadcastData := map[string]interface{}{"test": "broadcast"}
	jsonData, _ := json.Marshal(broadcastData)
	server.hub.broadcast <- &BroadcastMessage{
		Channel: "orderbook",
		Data:    jsonData,
	}
	
	// Wait for broadcast
	time.Sleep(200 * time.Millisecond)
	
	// Verify all clients received the broadcast
	for i, client := range clients {
		messages := client.GetReceivedMessages()
		assert.Greater(t, len(messages), 0, "Client %d didn't receive broadcast", i)
		client.Close()
	}
	
	assert.NoError(t, mock.ExpectationsWereMet())
}