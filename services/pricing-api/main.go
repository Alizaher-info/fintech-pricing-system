package main

import (
	"context"
	"log"
	"net"

	pricingv1 "pricingapi/gen/pricing/v1"

	"google.golang.org/grpc"
)

type server struct {
	// Ensure server implements the PricingServiceServer interface
	// This is done by embedding the UnimplementedPricingServiceServer struct
	// which provides default implementations for all methods.
	// This allows us to add new methods in the future without breaking existing clients.
	// This is a common pattern in gRPC to ensure forward compatibility.
	pricingv1.UnimplementedPricingServiceServer
}

func (s *server) Quote(ctx context.Context, r *pricingv1.QuoteRequest) (*pricingv1.QuoteResponse, error) {
	rate := 0.049
	if r.GetAmount() > 20000 {
		rate = 0.059
	}
	monthly := (r.GetAmount() * rate / 12) + (r.GetAmount() / float64(r.GetTermMonths()))
	return &pricingv1.QuoteResponse{
		InterestRate:   rate,
		Apr:            rate + 0.005,
		MonthlyPayment: monthly,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pricingv1.RegisterPricingServiceServer(s, &server{})
	log.Println("pricing-api gRPC listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
