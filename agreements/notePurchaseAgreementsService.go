package agreements

import (
	"context"
	"log"

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

func (service *NotePurchaseAgreementService) Save(ctx context.Context, agreement *interop.NotePurchaseAgreement) (*interop.NotePurchaseAgreement, error) {
	agreementModel, err := service.factory.NewPurchaseAgreementModel(ctx, agreement)

	if err != nil {
		return agreementModel.GetPayload(), err
	}

	result, err := agreementModel.Save()

	if err != nil {
		return nil, err
	}

	return result.GetPayload(), nil
}

func NewNotePurchaseAgreementService(factory interfaces.INotePurchaseAgreementModelFactory) interfaces.IAgreementsService {
	log.Printf("factory: %+#v \n\n", factory)
	return &NotePurchaseAgreementService{
		factory: factory,
	}
}
