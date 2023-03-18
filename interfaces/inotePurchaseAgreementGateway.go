package interfaces

import "context"

type INotePurchaseAgreementGateway interface {
	SaveNotePurchaseAgreement(context.Context, IAgreementModel) (IAgreementModel, error)
	Getall(ctx context.Context) ([]IAgreementModel, error)
	FindOne(ctx context.Context, model IAgreementModel) (IAgreementModel, error)
}
