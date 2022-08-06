package documents

import (
	"context"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type DocumentContentServiceGrpc struct {
	interfaces.IDocumentContentService
	notePurchaseAgreementServiceClient interop.NotePurchaseAgreementServiceClient
}

func (docs *DocumentContentServiceGrpc) BuildNotePurchaseAgreement(ctx context.Context, payload *interop.NotePurchaseAgreement) (interfaces.IDocumentContent, error) {
	documentResult, err := docs.notePurchaseAgreementServiceClient.GetNotePurchaseAgreementDoc(ctx, &interop.GetNotePurchaseAgreementDocRequest{
		Payload: payload,
	})

	if err != nil {
		return nil, err
	}

	return &DocumentContent{
		DocContent: documentResult.GetFileContent(),
	}, nil
}

func (docs *DocumentContentServiceGrpc) BuildAccountStatement() {

}

func NewDocumentContentManager(client interop.NotePurchaseAgreementServiceClient) interfaces.IDocumentContentService {

	return &DocumentContentServiceGrpc{
		notePurchaseAgreementServiceClient: client,
	}
}
