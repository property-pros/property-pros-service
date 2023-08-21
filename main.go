//go:generate go run ./build/gen_accessors.go -v
package main

// generate getters and setters with https://github.com/masaushi/accessory
import (
	ctx "context"
	"fmt"

	"github.com/google/wire"
	"github.com/vireocloud/property-pros-service/agreements"
	awss3 "github.com/vireocloud/property-pros-service/aws-s3"
	"github.com/vireocloud/property-pros-service/bootstrap"
	"github.com/vireocloud/property-pros-service/common"
	"github.com/vireocloud/property-pros-service/config"
	"github.com/vireocloud/property-pros-service/data"
	"github.com/vireocloud/property-pros-service/documents"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"github.com/vireocloud/property-pros-service/server/controllers"
	"github.com/vireocloud/property-pros-service/server/interceptors"
	"github.com/vireocloud/property-pros-service/users"
)

var UserSet wire.ProviderSet = wire.NewSet(
	data.NewUsersRepository,
	users.NewUserModel,
	NewUserModelFactory,
	users.NewUsersGateway,
	users.NewUsersService,
	controllers.NewAuthController)

var NotePuchaseAgreementSet wire.ProviderSet = wire.NewSet(
	data.NewGormDatabase,
	awss3.NewClient,
	data.NewAgreementsRepository,

	agreements.NewNotePurchaseAgreementModel,
	NewNotePurchaseAgreementModelFactory,
	agreements.NewNotePurchaseAgreementGateway,
	bootstrap.NewGrpcConnection,
	bootstrap.NewNotePurchaseAgreementClient,
	documents.NewDocumentContentService,
	agreements.NewNotePurchaseAgreementService,
	controllers.NewNotePurchaseAgreementController)

var StatementSet wire.ProviderSet = wire.NewSet(
	data.NewStatementsRepository,
	common.NewLogger,
	interop.NewStatementServiceClient,
	controllers.NewStatementController)

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

func Bootstrap() (*bootstrap.App, error) {

	wire.Build(
		config.NewConfig,
		UserSet,
		NotePuchaseAgreementSet,
		StatementSet,

		provideAuthenticationInterceptor,
		interceptors.NewConsumerDrivenContractTestingInterceptor,
		interceptors.NewController,
		interceptors.NewGrpcInterceptor,
		bootstrap.NewApp,
		documents.InitFileFixtures,
	)

	return nil, nil
}

func provideAuthenticationInterceptor(authService interfaces.IUsersService) (*interceptors.AuthValidationInterceptor, error) {
	return interceptors.NewAuthValidationInterceptor(authService, controllers.GRPC_AUTH_METHOD, controllers.GRPC_REGISTRATION_METHOD), nil
}

type Factory struct {
}

func (factory *Factory) NewPurchaseAgreementModel(context ctx.Context, agreement *interop.NotePurchaseAgreement) (interfaces.IAgreementModel, error) {
	// agreementModel, err := NotePurchaseAgreementInitializer()
	// fmt.Printf("NewPurchaseAgreementModel: %v\n", agreementModel)
	// if err != nil {
	// 	return nil, err
	// }

	// agreementModel.BaseModel = common.NewBaseModel[interop.NotePurchaseAgreement](agreement, context)
	// agreementModel.SetContext(context)
	// agreementModel.SetPayload(agreement)

	return &agreements.NotePurchaseAgreementModel{}, nil
}

func NewNotePurchaseAgreementModelFactory() interfaces.INotePurchaseAgreementModelFactory {
	return &Factory{}
}

// func NotePurchaseAgreementInitializer() (*agreements.NotePurchaseAgreementModel, error) {
// 	wire.Build(
// 		config.NewConfig,
// 		data.NewGormDatabase,
// 		UserSet,
// 		NotePuchaseAgreementSet)
// 	return nil, nil
// }

func (factory *Factory) NewUserModel(context ctx.Context, user *interop.User) (interfaces.IUserModel, error) {
	return &users.UserModel{}, nil
	// fmt.Println("here in NewUserModel")
	// userModel, err := UserModelInitializer()
	// if err != nil {
	// 	fmt.Printf("error occured: %v\n", err)
	// 	return nil, err
	// }
	// fmt.Println("model initialised")
	// fmt.Printf("userModel1: %v\n", userModel)

	// userModel.Context = context
	// userModel.Payload = user
	// fmt.Printf("userModel: %v\n", userModel)
	// return userModel, err
}

func NewUserModelFactory() interfaces.IUserModelFactory {
	return &Factory{}
}
