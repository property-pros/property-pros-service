package agreements

import (
	"context"

	"github.com/vireocloud/property-pros-service/data"
	"github.com/vireocloud/property-pros-service/interfaces"
)

type NotePurchaseAgreementGateway struct {
	repository interfaces.IRepository[data.NotePurchaseAgreement]
	factory    interfaces.INotePurchaseAgreementModelFactory
}

func NewNotePurchaseAgreementGateway(repository interfaces.IRepository[data.NotePurchaseAgreement], factory interfaces.INotePurchaseAgreementModelFactory) *NotePurchaseAgreementGateway {
	return &NotePurchaseAgreementGateway{
		repository: repository,
		factory:    factory,
	}
}

func (g *NotePurchaseAgreementGateway) SaveNotePurchaseAgreement(ctx context.Context, agreement data.NotePurchaseAgreement) (data.NotePurchaseAgreement, error) {
	_, err := g.repository.Save(&agreement)

	if err != nil {
		return data.NotePurchaseAgreement{}, err
	}

	return agreement, nil
}

func (g *NotePurchaseAgreementGateway) Getall(ctx context.Context) ([]interfaces.IAgreementModel, error) {
	return nil, nil
	// results := g.repository.Query(nil)
	// models := []interfaces.IAgreementModel{}

	// for _, result := range results {
	// 	model, err := g.factory.NewPurchaseAgreementModel(ctx, result)

	// 	if err != nil {
	// 		return models, err
	// 	}

	// 	models = append(models, model)
	// }

	// return models, nil
}

func (g *NotePurchaseAgreementGateway) FindOne(ctx context.Context, model data.NotePurchaseAgreement) (*data.NotePurchaseAgreement, error) {
	return g.repository.FindOne(&model)
}
