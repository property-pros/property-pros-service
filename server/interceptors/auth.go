package interceptors

import (
	"context"
	"fmt"

	"github.com/vireocloud/property-pros-service/interfaces"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthValidationInterceptor struct {
	authService        interfaces.IUsersService
	authMethod         string
	registrationMethod string
}

func NewAuthValidationInterceptor(authService interfaces.IUsersService, authMethod string, registrationMethod string) *AuthValidationInterceptor {
	return &AuthValidationInterceptor{
		authService,
		authMethod,
		registrationMethod,
	}
}

func (validator *AuthValidationInterceptor) Validate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if info.FullMethod != validator.registrationMethod && info.FullMethod != validator.authMethod {
		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("%v Failed;  Could not read metadata from context", info.FullMethod))
		}

		token := md.Get("authorization")[0]

		if validator.authService.IsValidToken(ctx, token) {
			return handler(ctx, req)
		}

		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("%v Failed;  Invalid auth token", info.FullMethod))
	}

	return handler(ctx, req)
}
