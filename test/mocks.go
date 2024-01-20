package test

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
	"google.golang.org/grpc/metadata"
)

type MockNotePurchaseAgreementService_GetNotePurchaseAgreementDocClient struct {
	interop.NotePurchaseAgreementService_GetNotePurchaseAgreementDocClient
}

func (c *MockNotePurchaseAgreementService_GetNotePurchaseAgreementDocClient) Recv() (*interop.GetNotePurchaseAgreementDocResponse, error) {
	// Return mock response
	return &interop.GetNotePurchaseAgreementDocResponse{
		FileContent: []byte("mock file content"),
	}, nil
}

func (c *MockNotePurchaseAgreementService_GetNotePurchaseAgreementDocClient) Header() (metadata.MD, error) {
	// Return mock headers
	return metadata.MD{}, nil
}

func (c *MockNotePurchaseAgreementService_GetNotePurchaseAgreementDocClient) Trailer() metadata.MD {
	// Return mock trailers
	return metadata.MD{}
}

func (c *MockNotePurchaseAgreementService_GetNotePurchaseAgreementDocClient) CloseSend() error {
	return nil
}

func (c *MockNotePurchaseAgreementService_GetNotePurchaseAgreementDocClient) Context() context.Context {
	return context.Background()
}

func (c *MockNotePurchaseAgreementService_GetNotePurchaseAgreementDocClient) SendMsg(m interface{}) error {
	return nil
}

func (c *MockNotePurchaseAgreementService_GetNotePurchaseAgreementDocClient) RecvMsg(m interface{}) error {
	return nil
}
