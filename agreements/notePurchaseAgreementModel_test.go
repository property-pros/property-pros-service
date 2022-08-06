package agreements

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

var testNotePurchaseAgreementModel *NotePurchaseAgreementModel
var baseModel interop.BaseModel[interop.NotePurchaseAgreement]
var documentContentService interfaces.IDocumentContentService
var notePurchaseAgreementGateway interfaces.INotePurchaseAgreementGateway
var userService interfaces.IUserService

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type NotePurchaseAgreementTestSuite struct {
	suite.Suite
	testNotePurchaseAgreementModel *NotePurchaseAgreementModel
	baseModel                      *interop.BaseModel[interop.NotePurchaseAgreement]
	documentContentService         *MockDocumentContentService
	testDocumentContent            *MockDocumentContent
	notePurchaseAgreementGateway   interfaces.INotePurchaseAgreementGateway
	userService                    interfaces.IUserService
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *NotePurchaseAgreementTestSuite) SetupTest() {
	suite.baseModel = &interop.BaseModel[interop.NotePurchaseAgreement]{
		Context: context.TODO(),
		Payload: new(interop.NotePurchaseAgreement),
	}
	suite.documentContentService = new(MockDocumentContentService)
	suite.testDocumentContent = new(MockDocumentContent)
	//  suite.notePurchaseAgreementGateway interfaces.INotePurchaseAgreementGateway
	//  suite.userService interfaces.IUserService
	suite.testNotePurchaseAgreementModel = &NotePurchaseAgreementModel{
		BaseModel:              suite.baseModel,
		documentContentService: suite.documentContentService}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestNotePurchaseAgreementTestSuite(t *testing.T) {
	suite.Run(t, new(NotePurchaseAgreementTestSuite))
}

func (suite *NotePurchaseAgreementTestSuite) TestGetDocumentContent() {
	suite.documentContentService.On("BuildNotePurchaseAgreement", suite.testNotePurchaseAgreementModel.Context, suite.testNotePurchaseAgreementModel.Payload).Return(suite.testDocumentContent, nil)
	suite.testDocumentContent.On("GetDocContent").Return([]byte{'a', 'b', 'c'})
	suite.testNotePurchaseAgreementModel.GetDocumentContent(nil)

	suite.documentContentService.AssertExpectations(suite.T())
	suite.testDocumentContent.AssertExpectations(suite.T())
}

func (suite *NotePurchaseAgreementTestSuite) TestSaveUser() {
	suite.documentContentService.On("BuildNotePurchaseAgreement", suite.testNotePurchaseAgreementModel.Context, suite.testNotePurchaseAgreementModel.Payload).Return(suite.testDocumentContent, nil)
	suite.testDocumentContent.On("GetDocContent").Return([]byte{'a', 'b', 'c'})
	suite.testNotePurchaseAgreementModel.SaveUser(nil)

	suite.documentContentService.AssertExpectations(suite.T())
	suite.testDocumentContent.AssertExpectations(suite.T())
}

type MockDocumentContentService struct {
	interfaces.IDocumentContentService
	mock.Mock
}

func (mocked *MockDocumentContentService) BuildNotePurchaseAgreement(ctx context.Context, payload *interop.NotePurchaseAgreement) (interfaces.IDocumentContent, error) {
	args := mocked.Called(ctx, payload)

	return args.Get(0).(interfaces.IDocumentContent), args.Error(1)
}

type MockDocumentContent struct {
	interfaces.IDocumentContent
	mock.Mock
}

func (mocked *MockDocumentContent) GetDocContent() []byte {
	args := mocked.Called()

	return args.Get(0).([]byte)
}

var _ interfaces.IDocumentContent = (*MockDocumentContent)(nil)
var _ interfaces.IDocumentContentService = (*MockDocumentContentService)(nil)

// func (suite *NotePurchaseAgreementTestSuite) TestSave(t *testing.T) {

// }

// func (suite *NotePurchaseAgreementTestSuite) TestSaveNotePurchaseAgreement(t *testing.T) {

// }

// func (suite *NotePurchaseAgreementTestSuite) TestGetUser(t *testing.T) {

// }

// func (suite *NotePurchaseAgreementTestSuite) TestNewNotePurchaseAgreementModel(t *testing.T) {

// }
