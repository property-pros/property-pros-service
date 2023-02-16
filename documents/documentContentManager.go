package documents

import (
	"context"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type DocumentContentService struct {
	interfaces.IDocumentContentService
	notePurchaseAgreementServiceDocClient interop.NotePurchaseAgreementServiceClient

	// documentGateway interfaces.IDocumentGateway
}

func (docs *DocumentContentService) BuildNotePurchaseAgreement(ctx context.Context, payload *interop.NotePurchaseAgreement) (interfaces.IDocumentContent, error) {
	documentResult, err := docs.notePurchaseAgreementServiceDocClient.GetNotePurchaseAgreementDoc(ctx, &interop.GetNotePurchaseAgreementDocRequest{
		Payload: payload,
	})

	if err != nil {
		return nil, err
	}

	return &DocumentContent{
		DocContent: documentResult.GetFileContent(),
	}, nil
}

func (docs *DocumentContentService) BuildAccountStatement() {

}

func NewDocumentContentManager(client interop.NotePurchaseAgreementServiceClient) interfaces.IDocumentContentService {

	return &DocumentContentService{
		notePurchaseAgreementServiceDocClient: client,
	}
}
