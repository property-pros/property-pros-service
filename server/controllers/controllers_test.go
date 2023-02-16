package controllers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type NotePurchaseAgreementsControllerTestSuite struct {
	suite.Suite
	controller            *NotePurchaseAgreementController
	testAgreementsService *MockTestAgreementsService
	testUserService       interfaces.IUsersService
}

func TestNotPurchaseAgreementsControllerSuite(t *testing.T) {
	suite.Run(t, new(NotePurchaseAgreementsControllerTestSuite))
}

func (suite *NotePurchaseAgreementsControllerTestSuite) TestGetNotePurchaseAgreemen() {
	t := suite.T()
	testDocContent := []byte("test doc")
	suite.testAgreementsService.On("GetNotePurchaseAgreementDocContent").Return(testDocContent, nil)

	result, err := suite.controller.GetNotePurchaseAgreementDoc(context.TODO(), &interop.GetNotePurchaseAgreementDocRequest{
		Payload: &interop.NotePurchaseAgreement{
			Id: "test",
		},
	})

	if err != nil {
		t.Error(err)
	}

	if result == nil {
		t.Errorf("Expected controller.GetNotePurchaseAgreementDoc to not return nil")
	}
}

func (suite *NotePurchaseAgreementsControllerTestSuite) SetupTest() {

	suite.testAgreementsService = &MockTestAgreementsService{}

	suite.controller = NewNotePurchaseAgreementController(suite.testAgreementsService, suite.testUserService)
}

type MockTestAgreementsService struct {
	mock.Mock
}

func (mock *MockTestAgreementsService) GetNotePurchaseAgreementDocContent(context.Context, interfaces.IModelPayload) ([]byte, error) {
	args := mock.Called()

	return args.Get(0).([]byte), args.Error(1)
}

func (mock *MockTestAgreementsService) GetNotePurchaseAgreement(context.Context, interfaces.IModelPayload) (*interop.NotePurchaseAgreement, error) {
	args := mock.Called()

	return args.Get(0).(*interop.NotePurchaseAgreement), args.Error(1)
}

func (mock *MockTestAgreementsService) GetNotePurchaseAgreements(context.Context) ([]interfaces.IAgreementModel, error) {
	args := mock.Called()

	return args.Get(0).([]interfaces.IAgreementModel), args.Error(1)
}

func (mock *MockTestAgreementsService) Save(context.Context, *interop.NotePurchaseAgreement) (*interop.NotePurchaseAgreement, error) {
	args := mock.Called()

	return args.Get(0).(*interop.NotePurchaseAgreement), args.Error(1)
}
