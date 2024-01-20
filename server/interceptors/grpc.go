package interceptors

import (
	"context"
	"fmt"

	"github.com/vireocloud/property-pros-service/interfaces"
	"google.golang.org/grpc"
)

type GrpcInterceptor struct {
	authService interfaces.IUsersService

	testRequestHandler           *ConsumerDrivenContractTestingInterceptor
	authValidationRequestHandler *AuthValidationInterceptor
}

func NewGrpcInterceptor(authService interfaces.IUsersService,
	testRequestHandler *ConsumerDrivenContractTestingInterceptor,
	authValidationRequestHandler *AuthValidationInterceptor) *GrpcInterceptor {

	return &GrpcInterceptor{
		authService,
		testRequestHandler,
		authValidationRequestHandler,
	}
}

func (i *GrpcInterceptor) HandleRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Println("HandleRequest started")
	testResult, err := i.testRequestHandler.Test(ctx, req, info, handler)

	if err != nil {
		// return handler(ctx, req)
		fmt.Println("HandleRequest auth validating")
		return i.authValidationRequestHandler.Validate(ctx, req, info, handler)
	}

	fmt.Printf("HandleRequest returning test result; Test err: %v", err)
	return testResult, err
}

