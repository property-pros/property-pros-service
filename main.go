//go:build wireinject
// +build wireinject

//go:generate go run ./build/gen_accessors.go -v

package main

// generate getters and setters with https://github.com/masaushi/accessory
import (

	// controllers "github.com/vireocloud/property-pros-service/server/controllers"
	// "github.com/vireocloud/property-pros-service/server/third_party"

	ctx "context"
	"fmt"
	"github.com/google/wire"
	"github.com/vireocloud/property-pros-service/agreements"
	"github.com/vireocloud/property-pros-service/bootstrap"
	"github.com/vireocloud/property-pros-service/config"
	"github.com/vireocloud/property-pros-service/data"
	"github.com/vireocloud/property-pros-service/documents"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"github.com/vireocloud/property-pros-service/server/controllers"
	"github.com/vireocloud/property-pros-service/server/interceptors"
	"github.com/vireocloud/property-pros-service/users"
	"log"
)

func main() {
	app, err := Bootstrap()
	fmt.Println("app bootstrap: ", app)
	if err != nil {
		panic(fmt.Errorf("failed to boostrap application: %w", err))
	}

	err = app.Run()

	if err != nil {
		panic(fmt.Errorf("failed to run application: %w", err))
	}

}

var UserSet wire.ProviderSet = wire.NewSet(
	data.NewUsersRepository,
	users.NewUsersGateway,
	users.NewUserModel,
	NewUserModelFactory,
	users.NewUsersService,
	controllers.NewAuthController)

var NotePuchaseAgreementSet wire.ProviderSet = wire.NewSet(
	data.NewAgreementsRepository,
	agreements.NewNotePurchaseAgreementGateway,
	agreements.NewNotePurchaseAgreementModel,
	NewNotePurchaseAgreementModelFactory,
	bootstrap.NewGrpcConnection,
	bootstrap.NewNotePurchaseAgreementClient,
	documents.NewDocumentContentManager,
	agreements.NewNotePurchaseAgreementService,
	controllers.NewNotePurchaseAgreementController)

func Bootstrap() (*bootstrap.App, error) {

	wire.Build(
		config.NewConfig,
		UserSet,
		NotePuchaseAgreementSet,
		provideAuthenticationInterceptor,
		interceptors.NewConsumerDrivenContractTestingInterceptor,
		interceptors.NewController,
		interceptors.NewGrpcInterceptor,
		bootstrap.NewApp,
	)

	return nil, nil
}

func provideAuthenticationInterceptor(authService interfaces.IUsersService) (*interceptors.AuthValidationInterceptor, error) {
	return interceptors.NewAuthValidationInterceptor(authService, controllers.GRPC_AUTH_METHOD, controllers.GRPC_REGISTRATION_METHOD), nil
}

type Factory struct {
}

func (factory *Factory) NewPurchaseAgreementModel(context ctx.Context, agreement *interop.NotePurchaseAgreement) (interfaces.IAgreementModel, error) {
	agreementModel, err := NotePurchaseAgreementInitializer()
	log.Printf("new agreement model: %+#v", agreementModel)
	if err != nil {
		return nil, err
	}

	agreementModel.Context = context
	agreementModel.Payload = agreement

	return agreementModel, err
}

func NewNotePurchaseAgreementModelFactory() interfaces.INotePurchaseAgreementModelFactory {
	return &Factory{}
}

func NotePurchaseAgreementInitializer() (*agreements.NotePurchaseAgreementModel, error) {
	wire.Build(
		config.NewConfig,
		data.NewGormDatabase,
		UserSet,
		NotePuchaseAgreementSet)
	return nil, nil
}

func (factory *Factory) NewUserModel(context ctx.Context, user *interop.User) (interfaces.IUserModel, error) {
	userModel, err := UserModelInitializer()
	userModel.Context = context
	userModel.Payload = user

	return userModel, err
}

func NewUserModelFactory() interfaces.IUserModelFactory {
	return &Factory{}
}

func UserModelInitializer() (*users.UserModel, error) {
	wire.Build(
		data.NewGormDatabase,
		UserSet,
	)
	return nil, nil
}
