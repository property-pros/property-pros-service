package interfaces

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
)

type IDocumentContentService interface {
	BuildNotePurchaseAgreement(ctx context.Context, payload *interop.NotePurchaseAgreement) (IDocumentContent, error)
}
