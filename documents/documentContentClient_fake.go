package documents

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
)

type ClientMock struct {
	interop.NotePurchaseAgreementServiceClient
}

func NewDocClientMock() interop.NotePurchaseAgreementServiceClient {
	return &ClientMock{}
}

func (*ClientMock) GetNotePurchaseAgreementDoc(ctx context.Context, in *interop.GetNotePurchaseAgreementDocRequest, opts ...interop.CallOption) (*interop.GetNotePurchaseAgreementDocResponse, error) {
	return &interop.GetNotePurchaseAgreementDocResponse{
		FileContent: []byte{1, 2, 3},
	}, nil
}
