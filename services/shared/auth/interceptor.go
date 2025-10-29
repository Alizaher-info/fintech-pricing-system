package auth

import (
	"context"
	"services/shared/jwt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager *jwt.Manager
}

func NewAuthInterceptor(jwtManager *jwt.Manager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager: jwtManager}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip auth for health checks
		if strings.HasSuffix(info.FullMethod, "Health") {
			return handler(ctx, req)
		}

		// Extract and validate token
		userCtx, err := a.authorize(ctx)
		if err != nil {
			return nil, err
		}

		// Add user context to request context
		return handler(userCtx, req)
	}
}

func (a *AuthInterceptor) authorize(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := a.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token: "+err.Error())
	}

	// Add user info to context
	userCtx := context.WithValue(ctx, "user_id", claims.UserID)
	userCtx = context.WithValue(userCtx, "user_email", claims.Email)
	userCtx = context.WithValue(userCtx, "user_role", claims.Role)

	return userCtx, nil
}

// Helper to extract user from context
func GetUserID(ctx context.Context) int64 {
	if userID, ok := ctx.Value("user_id").(int64); ok {
		return userID
	}
	return 0
}

func GetUserRole(ctx context.Context) string {
	if role, ok := ctx.Value("user_role").(string); ok {
		return role
	}
	return ""
}
