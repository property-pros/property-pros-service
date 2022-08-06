package interfaces

import (
	ctx "context"

	"github.com/vireocloud/property-pros-service/interop"
)

type INotePurchaseAgreementModelFactory interface {
	NewPurchaseAgreementModel(context ctx.Context, agreement *interop.NotePurchaseAgreement) (IAgreementModel, error)
}
