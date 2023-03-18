package documents

import (
	"context"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type DocumentContentService struct {
	notePurchaseAgreementServiceDocClient interop.NotePurchaseAgreementServiceClient
	uploader                              interfaces.IDocUploader
}

func (docs *DocumentContentService) BuildNotePurchaseAgreement(ctx context.Context, payload *interop.NotePurchaseAgreement) (interfaces.IDocumentContent, error) {
	documentResult, err := docs.notePurchaseAgreementServiceDocClient.GetNotePurchaseAgreementDoc(ctx, &interop.GetNotePurchaseAgreementDocRequest{
		Payload: payload,
	})

	if err != nil {
		return nil, err
	}

	return &DocumentContent{
		// DocContent: []byte("dummy result"),
		// TODO:Replace with actual result, once client is implemented
		DocContent: documentResult.GetFileContent(),
	}, nil
}

func (docService *DocumentContentService) CreateAndSaveNotePurchaseAgreementDoc(ctx context.Context, payload *interop.NotePurchaseAgreement) (string, error) {
	result, err := docService.BuildNotePurchaseAgreement(ctx, payload)
	if err != nil {
		return "", err
	}

	return docService.uploader.PutObject(ctx, result.GetDocContent())
}

func (docService *DocumentContentService) GetDocumentContent(ctx context.Context, id string) ([]byte, error) {
	return docService.uploader.GetObject(ctx, id)
}

func NewDocumentContentService(
	client interop.NotePurchaseAgreementServiceClient,
	uploader interfaces.IDocUploader) *DocumentContentService {

	return &DocumentContentService{
		notePurchaseAgreementServiceDocClient: client,
		uploader:                              uploader,
	}
}

// func (docs *DocumentContentService) BuildAccountStatement() {

// }
