package interceptors

import (
	"context"
	"fmt"
	"log"

	"github.com/vireocloud/property-pros-service/constants"
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
	fmt.Printf("full method: %v; registration method: %v; auth method: %v", info.FullMethod, validator.registrationMethod, validator.authMethod)
	if info.FullMethod != validator.registrationMethod && info.FullMethod != validator.authMethod {
		log.Printf("authenticating: %v", info.FullMethod)
		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("%v Failed;  Could not read metadata from context", info.FullMethod))
		}

		authMetadata := md.Get("authorization")

		log.Printf("auth metadata: %+v", authMetadata);

		token := authMetadata[0]
		userId := validator.authService.UserIdIfValidToken(ctx, token)
		if userId != "" {
			ctxWithUserId := context.WithValue(ctx, constants.UserIdKey, userId)
			return handler(ctxWithUserId, req)
		}

		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("%v Failed;  Invalid auth token", info.FullMethod))
	}

	return handler(ctx, req)
}
