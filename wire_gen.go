// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/vireocloud/property-pros-service/agreements"
	"github.com/vireocloud/property-pros-service/aws-s3"
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

// Injectors from main.go:

func Bootstrap() (*bootstrap.App, error) {
	iNotePurchaseAgreementModelFactory := NewNotePurchaseAgreementModelFactory()
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	db, err := data.NewGormDatabase(configConfig)
	if err != nil {
		return nil, err
	}
	iRepository := data.NewAgreementsRepository(db)
	interfacesIRepository := data.NewUsersRepository(db)
	notePurchaseAgreementGateway := agreements.NewNotePurchaseAgreementGateway(iRepository, interfacesIRepository, iNotePurchaseAgreementModelFactory)
	iUsersGateway := users.NewUsersGateway(interfacesIRepository)
	clientConnInterface := bootstrap.NewGrpcConnection(configConfig)
	notePurchaseAgreementServiceClient := bootstrap.NewNotePurchaseAgreementClient(clientConnInterface)
	iDocUploader := awss3.NewClient()
	documentContentService := documents.NewDocumentContentService(notePurchaseAgreementServiceClient, iDocUploader)
	iAgreementsService := agreements.NewNotePurchaseAgreementService(iNotePurchaseAgreementModelFactory, notePurchaseAgreementGateway, iUsersGateway, documentContentService)
	iUsersService := users.NewUsersService(iUsersGateway)
	notePurchaseAgreementController := controllers.NewNotePurchaseAgreementController(iAgreementsService, iUsersService)
	authController := controllers.NewAuthController(iAgreementsService, iUsersService)
	iStatementsRepository := data.NewStatementsRepository(db)
	statementController := controllers.NewStatementController(iStatementsRepository)
	propertyProsApiController := interceptors.NewController(iAgreementsService, iUsersService)
	consumerDrivenContractTestingInterceptor := interceptors.NewConsumerDrivenContractTestingInterceptor(propertyProsApiController)
	authValidationInterceptor, err := provideAuthenticationInterceptor(iUsersService)
	if err != nil {
		return nil, err
	}
	grpcInterceptor := interceptors.NewGrpcInterceptor(iUsersService, consumerDrivenContractTestingInterceptor, authValidationInterceptor)
	app := bootstrap.NewApp(notePurchaseAgreementController, authController, statementController, configConfig, grpcInterceptor, authValidationInterceptor)
	return app, nil
}

// main.go:

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

var UserSet wire.ProviderSet = wire.NewSet(data.NewUsersRepository, users.NewUserModel, NewUserModelFactory, users.NewUsersGateway, users.NewUsersService, controllers.NewAuthController)

var NotePuchaseAgreementSet wire.ProviderSet = wire.NewSet(data.NewGormDatabase, awss3.NewClient, data.NewAgreementsRepository, agreements.NewNotePurchaseAgreementModel, NewNotePurchaseAgreementModelFactory, agreements.NewNotePurchaseAgreementGateway, bootstrap.NewGrpcConnection, bootstrap.NewNotePurchaseAgreementClient, documents.NewDocumentContentService, agreements.NewNotePurchaseAgreementService, controllers.NewNotePurchaseAgreementController)

var StatementSet wire.ProviderSet = wire.NewSet(data.NewStatementsRepository, common.NewLogger, controllers.NewStatementController)

func provideAuthenticationInterceptor(authService interfaces.IUsersService) (*interceptors.AuthValidationInterceptor, error) {
	return interceptors.NewAuthValidationInterceptor(authService, controllers.GRPC_AUTH_METHOD, controllers.GRPC_REGISTRATION_METHOD), nil
}

type Factory struct {
}

func (factory *Factory) NewPurchaseAgreementModel(context2 context.Context, agreement *interop.NotePurchaseAgreement) (interfaces.IAgreementModel, error) {

	return &agreements.NotePurchaseAgreementModel{}, nil
}

func NewNotePurchaseAgreementModelFactory() interfaces.INotePurchaseAgreementModelFactory {
	return &Factory{}
}

func (factory *Factory) NewUserModel(context2 context.Context, user *interop.User) (interfaces.IUserModel, error) {
	return &users.UserModel{}, nil

}

func NewUserModelFactory() interfaces.IUserModelFactory {
	return &Factory{}
}
