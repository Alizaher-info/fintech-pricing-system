package websocket

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"services/pricing-fetcher/pkg/logger"

	"github.com/gorilla/websocket"
)

// BinanceClient represents a WebSocket client for Binance
type BinanceClient struct {
	conn           *websocket.Conn
	url            string
	symbols        []string
	priceChannel   chan PriceUpdate
	errorChannel   chan error
	reconnectDelay time.Duration
	mu             sync.RWMutex
	isConnected    bool
	droppedCount   int64 // Track dropped messages
	lastLogTime    time.Time // Last time we logged dropped messages
}

// PriceUpdate represents a price update from Binance WebSocket
type PriceUpdate struct {
	Symbol    string
	Price     float64
	Volume    float64
	Timestamp time.Time
}

// BinanceWSMessage represents the WebSocket message format from Binance
type BinanceWSMessage struct {
	EventType string `json:"e"` // Event type
	EventTime int64  `json:"E"` // Event time
	Symbol    string `json:"s"` // Symbol (e.g., "BTCUSDT")
	Price     string `json:"p"` // Price
	Quantity  string `json:"q"` // Quantity
	TradeTime int64  `json:"T"` // Trade time
}

// NewBinanceClient creates a new Binance WebSocket client
func NewBinanceClient(symbols []string) *BinanceClient {
	return &BinanceClient{
		url:            "wss://stream.binance.com:9443/stream",
		symbols:        symbols,
		priceChannel:   make(chan PriceUpdate, 50000), // Large buffer for high-frequency trades (BTC has 1000+ trades/4sec)
		errorChannel:   make(chan error, 10),
		reconnectDelay: 5 * time.Second,
		isConnected:    false,
	}
}

// PriceUpdates returns the channel for receiving price updates
func (c *BinanceClient) PriceUpdates() <-chan PriceUpdate {
	return c.priceChannel
}

// Errors returns the channel for receiving errors
func (c *BinanceClient) Errors() <-chan error {
	return c.errorChannel
}

// IsConnected returns the connection status
func (c *BinanceClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// setConnected sets the connection status
func (c *BinanceClient) setConnected(status bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.isConnected = status
}

// buildStreamURL builds the WebSocket URL for subscribing to multiple symbols
func (c *BinanceClient) buildStreamURL() string {
	if len(c.symbols) == 0 {
		return ""
	}

	// Build stream names (e.g., btcusdt@trade/ethusdt@trade/solusdt@trade)
	streams := ""
	for i, symbol := range c.symbols {
		// Convert BTC to btcusdt (Binance format - lowercase)
		streamName := fmt.Sprintf("%susdt@trade", strings.ToLower(symbol))
		if i > 0 {
			streams += "/"
		}
		streams += streamName
	}

	return fmt.Sprintf("%s?streams=%s", c.url, streams)
}

// Connect establishes WebSocket connection and starts listening
func (c *BinanceClient) Connect() error {
	wsURL := c.buildStreamURL()
	if wsURL == "" {
		return fmt.Errorf("no symbols provided for WebSocket connection")
	}

	logger.Info(fmt.Sprintf("Connecting to Binance WebSocket: %s", wsURL))

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to Binance WebSocket: %w", err)
	}

	c.conn = conn
	c.setConnected(true)
	logger.Info("Successfully connected to Binance WebSocket")

	// Start reading messages in a goroutine
	go c.readMessages()

	return nil
}

// readMessages continuously reads messages from WebSocket
func (c *BinanceClient) readMessages() {
	defer func() {
		c.setConnected(false)
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			logger.Error(fmt.Sprintf("WebSocket read error: %v", err))
			c.errorChannel <- err
			return
		}

		// Parse the message
		priceUpdate, err := c.parseMessage(message)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to parse message: %v", err))
			continue
		}

		// Send to price channel (non-blocking)
		select {
		case c.priceChannel <- priceUpdate:
			// Message sent successfully
		default:
			// Channel is full, drop message
			c.droppedCount++
			
			// Log dropped messages summary every 10 seconds (not every drop)
			now := time.Now()
			if now.Sub(c.lastLogTime) >= 10*time.Second {
				if c.droppedCount > 0 {
					logger.Error(fmt.Sprintf("⚠️ Dropped %d messages in last 10s (channel full, increase buffer or processing speed)", c.droppedCount))
					c.droppedCount = 0
				}
				c.lastLogTime = now
			}
		}
	}
}

// parseMessage parses Binance WebSocket message to PriceUpdate
func (c *BinanceClient) parseMessage(message []byte) (PriceUpdate, error) {
	var wsMsg struct {
		Data BinanceWSMessage `json:"data"`
	}

	if err := json.Unmarshal(message, &wsMsg); err != nil {
		return PriceUpdate{}, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// Parse price
	var price float64
	fmt.Sscanf(wsMsg.Data.Price, "%f", &price)

	// Parse quantity (volume)
	var quantity float64
	fmt.Sscanf(wsMsg.Data.Quantity, "%f", &quantity)

	// Extract symbol (remove USDT suffix: BTCUSDT -> BTC)
	symbol := wsMsg.Data.Symbol[:len(wsMsg.Data.Symbol)-4]

	return PriceUpdate{
		Symbol:    symbol,
		Price:     price,
		Volume:    quantity,
		Timestamp: time.UnixMilli(wsMsg.Data.TradeTime),
	}, nil
}

// Close closes the WebSocket connection
func (c *BinanceClient) Close() error {
	c.setConnected(false)
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
