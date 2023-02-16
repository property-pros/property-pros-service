package agreements

import (
	"context"

	"github.com/vireocloud/property-pros-service/interfaces"
)

type NotePurchaseAgreementGateway struct {
	repository interfaces.IAgreementsRepository
	factory    interfaces.INotePurchaseAgreementModelFactory
}

func NewNotePurchaseAgreementGateway(repository interfaces.IAgreementsRepository, factory interfaces.INotePurchaseAgreementModelFactory) interfaces.INotePurchaseAgreementGateway {
	return &NotePurchaseAgreementGateway{
		repository: repository,
		factory:    factory,
	}
}

func (g *NotePurchaseAgreementGateway) SaveNotePurchaseAgreement(ctx context.Context, model interfaces.IAgreementModel) (interfaces.IAgreementModel, error) {
	result, err := g.repository.Save(model.GetPayload())

	if err != nil {
		return nil, err
	}

	return g.factory.NewPurchaseAgreementModel(ctx, result)
}

func (g *NotePurchaseAgreementGateway) Getall(ctx context.Context) ([]interfaces.IAgreementModel, error) {
	results := g.repository.Query(nil)
	models := []interfaces.IAgreementModel{}

	for _, result := range results {
		model, err := g.factory.NewPurchaseAgreementModel(ctx, result)

		if err != nil {
			return models, err
		}

		models = append(models, model)
	}

	return models, nil
}

func (g *NotePurchaseAgreementGateway) FindOne(ctx context.Context, model interfaces.IAgreementModel) (interfaces.IAgreementModel, error) {
	result, err := g.repository.FindOne(model.GetPayload())

	if err != nil {
		return nil, err
	}

	return g.factory.NewPurchaseAgreementModel(ctx, result)
}
