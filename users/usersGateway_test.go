package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vireocloud/property-pros-service/data"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type UsersGatewayTestSuite struct {
	suite.Suite
	testUserPayload *interop.NotePurchaseAgreement
	mockUserRepo    *MockUserRepo
	usersGateway    interfaces.IUsersGateway
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *UsersGatewayTestSuite) SetupTest() {
	suite.testUserPayload = &interop.NotePurchaseAgreement{
		User: &interop.User{
			EmailAddress: "hjh",
			Password:     "dfsdf",
		},
		Id: "dsfds",
	}
	suite.mockUserRepo = new(MockUserRepo)
	suite.usersGateway = NewUsersGateway(suite.mockUserRepo)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUsersGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(UsersGatewayTestSuite))
}

func (suite *UsersGatewayTestSuite) TestSaveUser() {
	suite.mockUserRepo.On("Save").Return(&data.User{
		Id: "someID",
	})

	suite.usersGateway.SaveUser(context.TODO(), suite.testUserPayload)
	suite.AssertExpectationsSaveUser()
}

func (suite *UsersGatewayTestSuite) TestGetUserByUsername() {
	suite.mockUserRepo.On("Query").Return([]*data.User{{
		Id:           "someID",
		EmailAddress: "test@email.com",
	}})

	suite.usersGateway.GetUserByUsername("test@email.com")
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UsersGatewayTestSuite) AssertExpectationsSaveUser() {
	suite.mockUserRepo.AssertExpectations(suite.T())
}

type MockUserRepo struct {
	interfaces.IRepository[data.User]
	mock.Mock
	// waitGroup sync.WaitGroup
}

func (mocked *MockUserRepo) Save(user *data.User) (*data.User, error) {
	args := mocked.Called()

	return args.Get(0).(*data.User), nil
}

func (mocked *MockUserRepo) Query(user *data.User) []*data.User {
	args := mocked.Called()

	return args.Get(0).([]*data.User)
}

var _ interfaces.IRepository[data.User] = (*MockUserRepo)(nil)
