package main

import (

	// controllers "github.com/vireocloud/property-pros-service/server/controllers"
	// "github.com/vireocloud/property-pros-service/server/third_party"

	ctx "context"
	"fmt"

	"github.com/google/wire"
	"github.com/vireocloud/property-pros-service/agreements"
	"github.com/vireocloud/property-pros-service/bootstrap"
	"github.com/vireocloud/property-pros-service/config"
	"github.com/vireocloud/property-pros-service/documents"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"github.com/vireocloud/property-pros-service/server/controllers"
)

func main() {
	app, err := Bootstrap()

	if err != nil {
		panic(fmt.Errorf("failed to boostrap application: %w", err))
	}

	err = app.Run()

	if err != nil {
		panic(fmt.Errorf("failed to run application: %w", err))
	}

}

var NotePuchaseAgreementSet wire.ProviderSet = wire.NewSet(
	config.NewConfig,
	agreements.NewNotePurchaseAgreementModel,
	NewNotePurchaseAgreementModelFactory,
	bootstrap.NewGrpcConnection,
	bootstrap.NewNotePurchaseAgreementClient,
	documents.NewDocumentContentManager,
	agreements.NewNotePurchaseAgreementService,
	controllers.NewNotePurchaseAgreementController)

func Bootstrap() (*bootstrap.App, error) {

	wire.Build(
		NotePuchaseAgreementSet,
		bootstrap.NewApp,
	)

	return nil, nil
}

type Factory struct {
	interfaces.INotePurchaseAgreementModelFactory
}

func (factory *Factory) NewPurchaseAgreementModel(context ctx.Context, agreement *interop.NotePurchaseAgreement) (interfaces.IAgreementModel, error) {
	agreementModel, err := NotePurchaseAgreementInitializer()
	agreementModel.Context = context
	agreementModel.Payload = agreement

	return agreementModel, err
}

func NewNotePurchaseAgreementModelFactory() interfaces.INotePurchaseAgreementModelFactory {
	return &Factory{}
}

func NotePurchaseAgreementInitializer() (*agreements.NotePurchaseAgreementModel, error) {
	wire.Build(NotePuchaseAgreementSet)
	return nil, nil
}
