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
	baseModel       *common.BaseModel[interop.User]
	testUserPayload *interop.User
	userGateway     *MockUserGateway
	userService     interfaces.IUsersService
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *UsersServiceTestSuite) SetupTest() {
	suite.testUserPayload = new(interop.User)
	suite.userGateway = new(MockUserGateway)
	suite.userService = &UsersService{
		userGateway: suite.userGateway,
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUsersServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UsersServiceTestSuite))
}

func (suite *UsersServiceTestSuite) TestAuthenticateUser() {
	suite.userGateway.On("GetUserByUsername", "test@email.com").Return(&interop.NotePurchaseAgreement{
		User: &interop.User{
			Password:     "password",
			EmailAddress: "test@email.com",
			Id:           "someId",
		},
	})

	authResult, err := suite.userService.AuthenticateUser(context.TODO(), &interop.User{
		EmailAddress: "test@email.com",
		Password:     "password",
	})

	suite.AssertExpectationsAuthenticateUser(authResult, err)
}

func (suite *UsersServiceTestSuite) AssertExpectationsAuthenticateUser(result string, err error) {
	suite.Assert().Nil(err)
	suite.Assert().NotEmpty(result)
}

type MockUserGateway struct {
	mock.Mock
}

func (mocked *MockUserGateway) GetUserByUsername(userEmail string) (*interop.NotePurchaseAgreement, error) {
	args := mocked.Called(userEmail)

	return args.Get(0).(*interop.NotePurchaseAgreement), nil
}

func (mocked *MockUserGateway) SaveUser(ctx context.Context, user *interop.NotePurchaseAgreement) (*interop.NotePurchaseAgreement, error) {
	args := mocked.Called()

	return args.Get(0).(*interop.NotePurchaseAgreement), args.Error(1)
}

var _ interfaces.IUsersGateway = (*MockUserGateway)(nil)
