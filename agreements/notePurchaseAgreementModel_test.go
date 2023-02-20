package agreements

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vireocloud/property-pros-service/common"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type NotePurchaseAgreementTestSuite struct {
	suite.Suite
	baseModel                      *common.BaseModel[interop.NotePurchaseAgreement]
	testNotePurchaseAgreementModel *NotePurchaseAgreementModel
	testUserModel                  *MockUserModel
	testUserPayload                *interop.User
	documentContentService         *MockDocumentContentService
	testDocumentContent            *MockDocumentContent
	notePurchaseAgreementGateway   *MockNotePurchaseAgreementGateway
	userService                    *MockUserService
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *NotePurchaseAgreementTestSuite) SetupTest() {
	suite.baseModel = &common.BaseModel[interop.NotePurchaseAgreement]{
		Context: context.TODO(),
		Payload: new(interop.NotePurchaseAgreement),
	}

	suite.documentContentService = new(MockDocumentContentService)
	suite.testDocumentContent = new(MockDocumentContent)
	suite.notePurchaseAgreementGateway = new(MockNotePurchaseAgreementGateway)
	suite.userService = new(MockUserService)
	suite.testUserModel = &MockUserModel{}
	suite.testUserPayload = new(interop.User)

	suite.baseModel.Payload.User = suite.testUserPayload

	suite.testNotePurchaseAgreementModel = &NotePurchaseAgreementModel{
		BaseModel:                    suite.baseModel,
		notePurchaseAgreement:        suite.baseModel.Payload,
		documentContentService:       suite.documentContentService,
		userService:                  suite.userService,
		notePurchaseAgreementGateway: suite.notePurchaseAgreementGateway,
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestNotePurchaseAgreementTestSuite(t *testing.T) {
	suite.Run(t, new(NotePurchaseAgreementTestSuite))
}

func (suite *NotePurchaseAgreementTestSuite) TestSave() {
	suite.SetExpectationsSaveUser()
	suite.SetExpectationsSaveNotePurchaseAgreement()

	suite.testNotePurchaseAgreementModel.Save()

	suite.AssertExpectationsSaveUser()
	suite.AssertExpectationsSaveNotePurchaseAgreement()
}

func (suite *NotePurchaseAgreementTestSuite) TestGetDocumentContent() {
	suite.SetExpectationsGetDocumentContent()

	suite.testNotePurchaseAgreementModel.GetDocumentContent(nil)

	suite.AssertExpectationsGetDocumentContent()
}

func (suite *NotePurchaseAgreementTestSuite) TestSaveUser() {
	suite.SetExpectationsSaveUser()

	suite.testNotePurchaseAgreementModel.SaveUser(nil)

	suite.AssertExpectationsSaveUser()
}

func (suite *NotePurchaseAgreementTestSuite) TestSaveNotePurchaseAgreement() {
	suite.SetExpectationsSaveNotePurchaseAgreement()

	suite.testNotePurchaseAgreementModel.SaveNotePurchaseAgreement(nil)

	suite.AssertExpectationsSaveNotePurchaseAgreement()
}

func (suite *NotePurchaseAgreementTestSuite) SetExpectationsSaveNotePurchaseAgreement() {
	suite.notePurchaseAgreementGateway.waitGroup.Add(1)
	suite.notePurchaseAgreementGateway.On("SaveNotePurchaseAgreement").Return(suite.testNotePurchaseAgreementModel, nil)

}

func (suite *NotePurchaseAgreementTestSuite) SetExpectationsGetDocumentContent() {

	suite.documentContentService.On("BuildNotePurchaseAgreement", suite.testNotePurchaseAgreementModel.Context, suite.testNotePurchaseAgreementModel.Payload).Return(suite.testDocumentContent, nil)
	suite.testDocumentContent.On("GetDocContent").Return([]byte{'a', 'b', 'c'})
}

func (suite *NotePurchaseAgreementTestSuite) SetExpectationsSaveUser() {
	suite.userService.waitGroup.Add(1)
	suite.testUserModel.On("GetPayload").Return(suite.testUserPayload)

	suite.userService.On("SaveUser").Return(suite.testUserModel.GetPayload(), nil)
}

func (suite *NotePurchaseAgreementTestSuite) AssertExpectationsGetDocumentContent() {

	suite.documentContentService.AssertExpectations(suite.T())
}
func (suite *NotePurchaseAgreementTestSuite) AssertExpectationsSaveUser() {
	suite.userService.waitGroup.Wait()
	suite.testDocumentContent.AssertExpectations(suite.T())
	suite.userService.AssertExpectations(suite.T())
}

func (suite *NotePurchaseAgreementTestSuite) AssertExpectationsSaveNotePurchaseAgreement() {

	suite.notePurchaseAgreementGateway.waitGroup.Wait()
	suite.notePurchaseAgreementGateway.AssertExpectations(suite.T())
}

type MockUserModel struct {
	interfaces.IUserModel
	mock.Mock
	Payload *interop.User
}

func (mocked *MockUserModel) GetPayload() *interop.User {
	args := mocked.Called()

	return args.Get(0).(*interop.User)
}

type MockUserService struct {
	interfaces.IUsersService
	mock.Mock
	waitGroup sync.WaitGroup
}

func (mocked *MockUserService) SaveUser(ctx context.Context, model *interop.User) (*interop.User, error) {
	defer mocked.waitGroup.Done()
	args := mocked.Called()

	return args.Get(0).(*interop.User), args.Error(1)
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

type MockNotePurchaseAgreementGateway struct {
	interfaces.INotePurchaseAgreementGateway
	mock.Mock
	waitGroup sync.WaitGroup
}

func (mocked *MockNotePurchaseAgreementGateway) SaveNotePurchaseAgreement(context.Context, interfaces.IAgreementModel) (interfaces.IAgreementModel, error) {
	defer mocked.waitGroup.Done()
	args := mocked.Called()

	return args.Get(0).(interfaces.IAgreementModel), args.Error(1)
}

var _ interfaces.INotePurchaseAgreementGateway = (*MockNotePurchaseAgreementGateway)(nil)
var _ interfaces.IDocumentContent = (*MockDocumentContent)(nil)
var _ interfaces.IDocumentContentService = (*MockDocumentContentService)(nil)
var _ interfaces.IUsersService = (*MockUserService)(nil)
