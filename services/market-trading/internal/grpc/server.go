package grpc

import (
	"context"
	"fmt"
	"sync"

	"services/market-trading/internal/cache"
	"services/market-trading/internal/models"
	"services/market-trading/pkg/logger"

	pb "services/market-trading/gen/market/v1"
)

// Server implements the gRPC MarketDataService
type Server struct {
	pb.UnimplementedMarketDataServiceServer
	cache       *cache.RedisCache
	subscribers map[chan *models.PriceUpdate]struct{}
	mu          sync.RWMutex
}

// NewServer creates a new gRPC server
func NewServer(cache *cache.RedisCache) *Server {
	return &Server{
		cache:       cache,
		subscribers: make(map[chan *models.PriceUpdate]struct{}),
	}
}

// StreamPrices implements server-side streaming for real-time price updates
func (s *Server) StreamPrices(req *pb.SubscribeRequest, stream pb.MarketDataService_StreamPricesServer) error {
	logger.Infof("New client subscribed to price stream (symbols: %v)", req.Symbols)

	// Create a channel for this client
	priceChannel := make(chan *models.PriceUpdate, 100)

	// Register subscriber
	s.mu.Lock()
	s.subscribers[priceChannel] = struct{}{}
	s.mu.Unlock()

	// Unregister on exit
	defer func() {
		s.mu.Lock()
		delete(s.subscribers, priceChannel)
		close(priceChannel)
		s.mu.Unlock()
		logger.Info("Client unsubscribed from price stream")
	}()

	// Send initial prices from cache if requested
	if len(req.Symbols) > 0 {
		ctx := context.Background()
		for _, symbol := range req.Symbols {
			price, err := s.cache.GetPrice(ctx, symbol)
			if err != nil {
				logger.Warnf("Failed to get initial price for %s: %v", symbol, err)
				continue
			}
			if price != nil {
				if err := stream.Send(convertPriceToProto(price)); err != nil {
					return err
				}
			}
		}
	}

	// Stream prices as they come in
	for {
		select {
		case price := <-priceChannel:
			// Filter by symbols if specified
			if len(req.Symbols) > 0 && !contains(req.Symbols, price.Symbol) {
				continue
			}

			// Send price update to client
			protoPrice := convertPriceToProto(price)
			if err := stream.Send(protoPrice); err != nil {
				logger.Errorf("Failed to send price update: %v", err)
				return err
			}

		case <-stream.Context().Done():
			logger.Info("Client disconnected from stream")
			return nil
		}
	}
}

// GetPrice implements unary RPC to get current price from cache
func (s *Server) GetPrice(ctx context.Context, req *pb.GetPriceRequest) (*pb.GetPriceResponse, error) {
	if req.Symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	// Get price from Redis cache
	price, err := s.cache.GetPrice(ctx, req.Symbol)
	if err != nil {
		logger.Errorf("Failed to get price for %s: %v", req.Symbol, err)
		return nil, fmt.Errorf("failed to get price: %w", err)
	}

	if price == nil {
		return &pb.GetPriceResponse{
			Price: nil,
			Found: false,
		}, nil
	}

	return &pb.GetPriceResponse{
		Price: convertPriceToProto(price),
		Found: true,
	}, nil
}

// BroadcastPrice broadcasts a price update to all connected clients
func (s *Server) BroadcastPrice(price *models.PriceUpdate) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Send to all subscribers (non-blocking)
	for ch := range s.subscribers {
		select {
		case ch <- price:
			// Successfully sent
		default:
			// Channel full, skip this update to avoid blocking
			logger.Warnf("Subscriber channel full, skipping update for %s", price.Symbol)
		}
	}
}

// GetSubscriberCount returns the number of active subscribers
func (s *Server) GetSubscriberCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.subscribers)
}

// convertPriceToProto converts internal price model to protobuf message
func convertPriceToProto(price *models.PriceUpdate) *pb.PriceUpdate {
	return &pb.PriceUpdate{
		Symbol:    price.Symbol,
		Price:     price.Price,
		BidPrice:  price.BidPrice,
		AskPrice:  price.AskPrice,
		Volume_24H: price.Volume24h,
		Change_24H: price.Change24h,
		High_24H:   price.High24h,
		Low_24H:    price.Low24h,
		Timestamp: price.Timestamp.UnixMilli(),
	}
}

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
