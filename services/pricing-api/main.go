package main

import (
	"context"
	"log"
	"net"
	"os"

	pricingv1 "services/pricing-api/gen/pricing/v1"
	"services/shared/auth"
	"services/shared/jwt"

	"google.golang.org/grpc"
)

type server struct {
	pricingv1.UnimplementedPricingServiceServer
}

func (s *server) Quote(ctx context.Context, r *pricingv1.QuoteRequest) (*pricingv1.QuoteResponse, error) {
	// Extract user context from auth interceptor
	userID := auth.GetUserID(ctx)
	userRole := auth.GetUserRole(ctx)

	log.Printf("Processing quote for user %d (role: %s), amount: %.2f",
		userID, userRole, r.GetAmount())

	// Base rate calculation
	rate := 0.049

	// Premium users get better rates
	if userRole == "premium" {
		rate = 0.039
	}

	// Large amounts get different rates
	if r.GetAmount() > 20000 {
		rate += 0.01
	}

	// Risk score affects rate
	if r.GetRiskScore() > 0 {
		if r.GetRiskScore() < 600 {
			rate += 0.02 // Higher risk, higher rate
		} else if r.GetRiskScore() > 750 {
			rate -= 0.005 // Lower risk, lower rate
		}
	}

	monthly := (r.GetAmount() * rate / 12) + (r.GetAmount() / float64(r.GetTermMonths()))

	return &pricingv1.QuoteResponse{
		InterestRate:   rate,
		Apr:            rate + 0.005,
		MonthlyPayment: monthly,
	}, nil
}

func main() {
	// JWT configuration
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // Default for development
	}

	jwtManager := jwt.NewManager(jwtSecret, "fintech-api")
	authInterceptor := auth.NewAuthInterceptor(jwtManager)

	// Create gRPC server with auth interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	pricingv1.RegisterPricingServiceServer(s, &server{})
	log.Println("Secure pricing-api gRPC listening on :50051")

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
