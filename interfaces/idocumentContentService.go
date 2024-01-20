package interfaces

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
)

type IDocumentContentService interface {
	BuildNotePurchaseAgreement(ctx context.Context, payload *interop.NotePurchaseAgreement) (IDocumentContent, error)
	BuildStatement(ctx context.Context, payload *interop.Statement) (IDocumentContent, error)
	GetDocumentContent(ctx context.Context, docURL string) ([]byte, error)
	CreateAndSaveNotePurchaseAgreementDoc(ctx context.Context, payload *interop.NotePurchaseAgreement) (string, error)
}
