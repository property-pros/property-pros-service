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
type UsersGatewayTestSuite struct {
	suite.Suite
	baseModel            *common.BaseModel[interop.User]
	mockUserModelFactory *MockUserModelFactory
	mockUserModel        *MockUserModel
	testUserPayload      *interop.User
	mockUserRepo         *MockUserRepo
	usersGateway         interfaces.IUsersGateway
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *UsersGatewayTestSuite) SetupTest() {
	suite.testUserPayload = new(interop.User)
	suite.baseModel = &common.BaseModel[interop.User]{
		Context: context.TODO(),
		Payload: new(interop.User),
	}
	suite.mockUserModelFactory = new(MockUserModelFactory)
	suite.mockUserRepo = new(MockUserRepo)
	suite.usersGateway = &UsersGateway{
		factory: suite.mockUserModelFactory,
		repo:    suite.mockUserRepo,
	}
	suite.mockUserModel = &MockUserModel{
		BaseModel: suite.baseModel,
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUsersGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(UsersGatewayTestSuite))
}

func (suite *UsersGatewayTestSuite) TestSaveUser() {
	suite.SetExpectationsSaveUser()

	suite.usersGateway.SaveUser(suite.mockUserModel)

	suite.AssertExpectationsSaveUser()
}

func (suite *UsersGatewayTestSuite) SetExpectationsSaveUser() {
	// suite.userService.waitGroup.Add(1)
	suite.mockUserModelFactory.On("NewUserModel").Return(suite.mockUserModel, nil)
	suite.mockUserModel.On("GetPayload").Return(suite.testUserPayload, nil)
	suite.mockUserRepo.On("Save").Return(suite.testUserPayload, nil)
}

func (suite *UsersGatewayTestSuite) AssertExpectationsSaveUser() {
	// suite.userService.waitGroup.Wait()
	suite.mockUserModelFactory.AssertExpectations(suite.T())
	suite.mockUserModel.AssertExpectations(suite.T())
	suite.mockUserRepo.AssertExpectations(suite.T())
}

type MockUserRepo struct {
	interfaces.IRepository[interop.User]
	mock.Mock
	// waitGroup sync.WaitGroup
}

func (mocked *MockUserRepo) Save(user *interop.User) (*interop.User, error) {
	// defer mocked.waitGroup.Done()
	args := mocked.Called()

	return args.Get(0).(*interop.User), args.Error(1)
}

var _ interfaces.IRepository[interop.User] = (*MockUserRepo)(nil)

var _ interfaces.IUserModel = (*MockUserModel)(nil)
