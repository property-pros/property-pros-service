package agreements

import (
	"context"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type NotePurchaseAgreementService struct {
	interfaces.IAgreementsService
	factory interfaces.INotePurchaseAgreementModelFactory
}

func (service *NotePurchaseAgreementService) GetNotePurchaseAgreementDocContent(context.Context, *interop.NotePurchaseAgreement) ([]byte, error) {
	return nil, nil
}

func NewNotePurchaseAgreementService(factory interfaces.INotePurchaseAgreementModelFactory) interfaces.IAgreementsService {
	return &NotePurchaseAgreementService{
		factory: factory,
	}
}
