package interfaces

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
)

type IAgreementsService interface {
	GetNotePurchaseAgreements(ctx context.Context, userID string) ([]*interop.NotePurchaseAgreement, error)
	GetNotePurchaseAgreement(context.Context, IModelPayload) (*interop.NotePurchaseAgreement, error)
	GetNotePurchaseAgreementDocContent(context.Context, IModelPayload) ([]byte, error)
	Save(context.Context, *interop.NotePurchaseAgreement) (*interop.NotePurchaseAgreement, error)
}
