package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vireocloud/property-pros-service/common"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type UsersServiceTestSuite struct {
	suite.Suite
	baseModel            *common.BaseModel[interop.User]
	mockUserModelFactory *MockUserModelFactory
	testUserModel        *MockUserModel
	testUserPayload      *interop.User
	userGateway          *MockUserGateway
	userService          interfaces.IUsersService
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *UsersServiceTestSuite) SetupTest() {
	suite.testUserPayload = new(interop.User)
	suite.baseModel = &common.BaseModel[interop.User]{
		Context: context.TODO(),
		Payload: new(interop.User),
	}
	suite.userGateway = new(MockUserGateway)
	suite.mockUserModelFactory = new(MockUserModelFactory)
	suite.userService = &UsersService{
		factory: suite.mockUserModelFactory,
	}
	suite.testUserModel = &MockUserModel{
		BaseModel: suite.baseModel,
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUsersServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UsersServiceTestSuite))
}

func (suite *UsersServiceTestSuite) TestSaveUser() {
	suite.SetExpectationsSaveUser()

	suite.userService.SaveUser(context.TODO(), suite.testUserPayload)

	suite.AssertExpectationsSaveUser()
}

func (suite *UsersServiceTestSuite) TestAuthenticateUser() {
	suite.SetExpectationsAuthenticateUser()

	authResult, err := suite.userService.AuthenticateUser(context.TODO(), suite.testUserPayload)

	suite.AssertExpectationsAuthenticateUser(authResult, err)
}

func (suite *UsersServiceTestSuite) SetExpectationsAuthenticateUser() {
	// suite.userService.waitGroup.Add(1)
	suite.mockUserModelFactory.On("NewUserModel").Return(suite.testUserModel, nil)
	suite.testUserModel.On("HasAuthenticIdentity").Return(true, nil)
	suite.testUserModel.On("HasAuthorization").Return(true, nil)

}

func (suite *UsersServiceTestSuite) AssertExpectationsAuthenticateUser(result bool, err error) {

	suite.Assert().Nil(err)
	suite.Assert().True(result)
	// suite.userService.waitGroup.Wait()
	suite.mockUserModelFactory.AssertExpectations(suite.T())
	suite.testUserModel.AssertExpectations(suite.T())
}
func (suite *UsersServiceTestSuite) SetExpectationsSaveUser() {
	// suite.userService.waitGroup.Add(1)
	suite.mockUserModelFactory.On("NewUserModel").Return(suite.testUserModel, nil)
	suite.testUserModel.On("Save").Return(suite.testUserModel, nil)
	suite.testUserModel.On("GetPayload").Return(suite.testUserPayload, nil)
}

func (suite *UsersServiceTestSuite) AssertExpectationsSaveUser() {
	// suite.userService.waitGroup.Wait()
	suite.mockUserModelFactory.AssertExpectations(suite.T())
	suite.testUserModel.AssertExpectations(suite.T())
}

type MockUserModelFactory struct {
	mock.Mock
}

func (mocked *MockUserModelFactory) NewUserModel(context.Context, *interop.User) (interfaces.IUserModel, error) {
	args := mocked.Called()

	return args.Get(0).(interfaces.IUserModel), args.Error(1)
}

type MockUserGateway struct {
	mock.Mock
}

func (mocked *MockUserGateway) GetUserByUsername(interfaces.IUserModel) (interfaces.IUserModel, error) {
	args := mocked.Called()

	return args.Get(0).(interfaces.IUserModel), args.Error(1)
}

var _ interfaces.IUsersGateway = (*MockUserGateway)(nil)

type MockUserModel struct {
	mock.Mock
	*common.BaseModel[interop.User]
}

func (mocked *MockUserModel) HasAuthenticIdentity() (bool, error) {

	mocked.Called()

	return true, nil
}
func (mocked *MockUserModel) HasAuthorization() (bool, error) {

	mocked.Called()

	return true, nil
}

func (mocked *MockUserModel) Save() (interfaces.IUserModel, error) {

	mocked.Called()

	return mocked, nil
}

func (mocked *MockUserModel) MatchCredentials(identity interfaces.IUserModel) (bool, error) {
	mocked.Called()

	return true, nil
}

func (mocked *MockUserModel) GetPayload() *interop.User {

	mocked.Called()

	return mocked.BaseModel.Payload
}

func (mocked *MockUserModel) GetPassword() string {

	mocked.Called()

	return mocked.BaseModel.Payload.GetPassword()
}

func (mocked *MockUserModel) GetId() string {

	mocked.Called()

	return mocked.BaseModel.Payload.GetId()
}

var _ interfaces.IUserModel = (*MockUserModel)(nil)
