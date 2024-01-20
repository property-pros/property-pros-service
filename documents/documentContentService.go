package documents

import (
	"context"
	"fmt"
	"io"

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
	
	fmt.Printf("building note purchase agreement doc content; payload filecontent length: %v\n", len(payload.FileContent))
	
	documentStream, err := docs.notePurchaseAgreementServiceDocClient.GetNotePurchaseAgreementDoc(ctx, &interop.GetNotePurchaseAgreementDocRequest{
		Payload: payload,
	})
	
	fmt.Println("got docstream from notepurchase agreement service")
	if err != nil {
		fmt.Printf("error getting note purchase agreement doc from doc service: %s\n", err.Error())
		return nil, err
	}

	var documentResult interop.GetNotePurchaseAgreementDocResponse

	fmt.Println("reading chunks from stream")
	for {
		chunk, err := documentStream.Recv(); 
		if err == io.EOF {
			break;
		}
		fmt.Println("chunk err: ", err)
		fmt.Println("got chunk: ", chunk.GetFileContent())
		documentResult.FileContent = append(documentResult.FileContent, chunk.GetFileContent()...)
	}



	fmt.Println("done reading chunks from stream")
	return &DocumentContent{
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
	fmt.Println("building note purchase agreement doc")
	result, err := docService.BuildNotePurchaseAgreement(ctx, payload)
	if err != nil {
		return "", err
	}
	fmt.Println("built note purchase agreement doc")
	return docService.uploader.PutObject(ctx, result.GetDocContent())
}

func (docService *DocumentContentService) GetDocumentContent(ctx context.Context, id string) ([]byte, error) {
	return docService.uploader.GetObject(ctx, id)
}

// func (docs *DocumentContentService) BuildAccountStatement() {

// }
