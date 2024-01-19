package interceptors

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ConsumerDrivenContractTestingInterceptor struct {
	controller *PropertyProsApiController
}

func NewConsumerDrivenContractTestingInterceptor(controller *PropertyProsApiController) *ConsumerDrivenContractTestingInterceptor {
	return &ConsumerDrivenContractTestingInterceptor{
		controller: controller,
	}
}

func (interceptor *ConsumerDrivenContractTestingInterceptor) Test(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (result interface{}, err error) {

	fmt.Println("full method: ", info.FullMethod)

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, fmt.Errorf("%v Failed;  Could not read metadata from context", info.FullMethod)
	}

	isTestMetadataCollection := md.Get("is-test")

	if isTestMetadataCollection != nil {

		isTestMetadata := md.Get("is-test")[0]

		isTest, err := strconv.ParseBool(isTestMetadata)

		if err == nil {

			if isTest {
				functionName := cleanMethod(info.FullMethod)

				testServerReflectedValue := reflect.ValueOf(interceptor.controller)

				reflectedMethod := testServerReflectedValue.MethodByName(functionName)

				reflectedParamValues := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)}

				results := reflectedMethod.Call(reflectedParamValues)

				errorResult := results[1].Interface()
				err = nil

				if errorResult != nil {
					err = errorResult.(error)
				}

				return results[0].Interface(), err
			}
		}
	}

	return result, fmt.Errorf("this is not a test")
}

func cleanMethod(grpcMethodName string) string {

	return strings.Split(grpcMethodName, "/")[2]
}

type PropertyProsApiController struct {
	interop.UnsafeAuthenticationServiceServer
	interop.UnsafeNotePurchaseAgreementServiceServer

	authService                  interfaces.IUsersService
	notePurchaseAgreementService interfaces.IAgreementsService
}

func (c *PropertyProsApiController) SaveNotePurchaseAgreement(ctx context.Context, req *interop.SaveNotePurchaseAgreementRequest) (response *interop.SaveNotePurchaseAgreementResponse, errResult error) {
	return &interop.SaveNotePurchaseAgreementResponse{
		Payload: &interop.RecordResultPayload{
			Id: "36c43851-6caa-4c48-962e-0735a6e27e67",
		},
	}, nil
}

func (c *PropertyProsApiController) GetNotePurchaseAgreementDoc(*interop.GetNotePurchaseAgreementDocRequest, interop.NotePurchaseAgreementService_GetNotePurchaseAgreementDocServer) error {
	return nil
}

func (c *PropertyProsApiController) GetNotePurchaseAgreement(ctx context.Context, req *interop.GetNotePurchaseAgreementRequest) (response *interop.GetNotePurchaseAgreementResponse, errResult error) {
	return &interop.GetNotePurchaseAgreementResponse{
			Payload: &interop.NotePurchaseAgreement{
				Id:          uuid.New().String(),
				FileContent: []byte("test file content"),
			},
		},
		nil
}

func (c *PropertyProsApiController) GetNotePurchaseAgreements(ctx context.Context, req *interop.GetNotePurchaseAgreementsRequest) (response *interop.GetNotePurchaseAgreementsResponse, errResult error) {
	return &interop.GetNotePurchaseAgreementsResponse{
			Payload: []*interop.RecordResultPayload{
				{Id: uuid.New().String()},
			},
		},

		nil
}

func (c *PropertyProsApiController) AuthenticateUser(ctx context.Context, req *interop.AuthenticateUserRequest) (*interop.AuthenticateUserResponse, error) {

	response := &interop.AuthenticateUserResponse{}

	response.IsAuthenticated = true

	// We want to extract metadata from the incomming context.
	// We dont create a new context since we dont wanna overwrite old metadata
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("could not grab metadata from context")
	}
	// Set authorization  metadata for the client to send in subsequent requests
	meta.Set("authorization", c.authService.GenerateBasicUserAuthToken(req.Payload))
	// Metadata is sent on its own, so we need to send the header. There is also something called Trailer
	grpc.SendHeader(ctx, meta)

	return response, nil
}

func NewController(notePurchaseAgreementService interfaces.IAgreementsService, authService interfaces.IUsersService) *PropertyProsApiController {
	return &PropertyProsApiController{
		notePurchaseAgreementService: notePurchaseAgreementService,
		authService:                  authService,
	}
}

var _ interop.NotePurchaseAgreementServiceServer = (*PropertyProsApiController)(nil)

var _ interop.AuthenticationServiceServer = (*PropertyProsApiController)(nil)
