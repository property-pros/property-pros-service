package documents

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
	"github.com/vireocloud/property-pros-service/test"
	"google.golang.org/grpc"
)

type ClientMock struct {
	interop.NotePurchaseAgreementServiceClient
}

func NewDocClientMock() interop.NotePurchaseAgreementServiceClient {
	return &ClientMock{}
}

func (*ClientMock) GetNotePurchaseAgreementDoc(ctx context.Context, in *interop.GetNotePurchaseAgreementDocRequest, opts ...grpc.CallOption) (interop.NotePurchaseAgreementService_GetNotePurchaseAgreementDocClient, error) {
	return &test.MockNotePurchaseAgreementService_GetNotePurchaseAgreementDocClient{}, nil
}