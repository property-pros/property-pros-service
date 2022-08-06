package interfaces

import "context"

type INotePurchaseAgreementGateway interface {
	SaveNotePurchaseAgreement(context.Context, IAgreementModel) (IAgreementModel, error)
}
