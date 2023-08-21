package documents

import (
	"context"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type DocumentContentService struct {
	notePurchaseAgreementServiceDocClient interop.NotePurchaseAgreementServiceClient
	statementServiceClient                interop.StatementServiceClient
	uploader                              interfaces.IDocUploader
}

func NewDocumentContentService(
	npaClient interop.NotePurchaseAgreementServiceClient,
	statementClient interop.StatementServiceClient,
	uploader interfaces.IDocUploader) interfaces.IDocumentContentService {

	return &DocumentContentService{
		notePurchaseAgreementServiceDocClient: npaClient,
		statementServiceClient:                statementClient,
		uploader:                              uploader,
	}
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

func (docs *DocumentContentService) BuildStatement(ctx context.Context, payload *interop.Statement) (interfaces.IDocumentContent, error) {
	documentResult, err := docs.statementServiceClient.GetStatementDoc(ctx, &interop.GetStatementDocRequest{
		Payload: payload,
	})

	if err != nil {
		return nil, err
	}

	return &DocumentContent{
		DocContent: documentResult.GetDocument(),
	}, err
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

// func (docs *DocumentContentService) BuildAccountStatement() {

// }
