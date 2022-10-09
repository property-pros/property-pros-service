package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vireocloud/property-pros-service/common"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type UserModelTestSuite struct {
	suite.Suite
	baseModel       *common.BaseModel[interop.User]
	testUserPayload *interop.User
	userGateway     *MockUserGateway
	userModel       interfaces.IUserModel
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *UserModelTestSuite) SetupTest() {
	suite.testUserPayload = new(interop.User)
	suite.baseModel = &common.BaseModel[interop.User]{
		Context: context.TODO(),
		Payload: suite.testUserPayload,
	}
	suite.userGateway = new(MockUserGateway)

	suite.userModel = &UserModel{
		User:      suite.testUserPayload,
		BaseModel: suite.baseModel,
		gateway:   suite.userGateway,
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserModelTestSuite(t *testing.T) {
	suite.Run(t, new(UserModelTestSuite))
}

func (suite *UserModelTestSuite) TestSave() {
	suite.SetExpectationsSaveUser()

	suite.userModel.Save()

	suite.AssertExpectationsSaveUser()
}

func (suite *UserModelTestSuite) TestHasAuthenticIdentity() {
	suite.SetExpectationsHasAuthenticIdentity()

	isIdentityAuthentic, authenticationError := suite.userModel.HasAuthenticIdentity()

	suite.AssertExpectationsHasAuthenticIdentity(isIdentityAuthentic, authenticationError)
}

func (suite *UserModelTestSuite) TestHasAuthorization() {

	suite.SetExpectationsHasAuthorization()

	isIdentityAuthentic, authenticationError := suite.userModel.HasAuthorization()

	suite.AssertExpectationsHasAuthorization(isIdentityAuthentic, authenticationError)
}

func (suite *UserModelTestSuite) SetExpectationsHasAuthorization() {
	suite.userModel.(*UserModel).Id = "testid"
}

func (suite *UserModelTestSuite) AssertExpectationsHasAuthorization(isAuthorized bool, authorizationError error) {
	suite.Assert().True(isAuthorized)
	suite.Assert().Nil(authorizationError)
}

func (suite *UserModelTestSuite) SetExpectationsHasAuthenticIdentity() {
	suite.userGateway.On("GetUserByUsername").Return(suite.userModel, nil)
}

func (suite *UserModelTestSuite) AssertExpectationsHasAuthenticIdentity(isIdentityAuthentic bool, authenticationError error) {
	suite.Assert().True(isIdentityAuthentic)
	suite.Assert().Nil(authenticationError)
	suite.userGateway.AssertExpectations(suite.T())
}

func (suite *UserModelTestSuite) SetExpectationsSaveUser() {
	suite.userGateway.On("SaveUser").Return(suite.userModel, nil)
}

func (suite *UserModelTestSuite) AssertExpectationsSaveUser() {
	suite.userGateway.AssertExpectations(suite.T())
}

var _ interfaces.IUserModel = (*MockUserModel)(nil)

func (mocked *MockUserGateway) SaveUser(user interfaces.IUserModel) (interfaces.IUserModel, error) {
	args := mocked.Called()

	return args.Get(0).(*UserModel), args.Error(1)
}
