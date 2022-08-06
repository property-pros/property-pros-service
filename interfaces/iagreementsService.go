package interfaces

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
)

type IAgreementsService interface {
	GetNotePurchaseAgreementDocContent(context.Context, *interop.NotePurchaseAgreement) ([]byte, error)
}
