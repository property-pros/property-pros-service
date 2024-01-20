package controllers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type StatementControllerTestSuite struct {
	suite.Suite
	controller        *StatementController
	testStatementRepo *MockStatementRepo
}

func (suite *StatementControllerTestSuite) SetupTest() {
	suite.testStatementRepo = &MockStatementRepo{}
	suite.controller = NewStatementController(suite.testStatementRepo)
}
func (suite *StatementControllerTestSuite) TestGetStatements() {
	ctx := context.TODO()
	userId := "someUserId"
	expectedStatements := []*interop.Statement{
		{
			Id: "statementId1",
		},
		{
			Id: "statementId2",
		},
	}

	suite.testStatementRepo.On("Query", mock.IsType(&interop.Statement{
		UserId: userId,
	})).Return(expectedStatements)

	result, err := suite.controller.GetStatements(ctx, &interop.GetStatementsRequest{
		UserId: userId,
	})

	suite.Nil(err)
	suite.Equal(expectedStatements, result.Payload.Statements)
	suite.testStatementRepo.AssertExpectations(suite.T())
}

func TestStatementController(t *testing.T) {
	suite.Run(t, new(StatementControllerTestSuite))
}


type MockStatementRepo struct {
	mock.Mock
}

func (m *MockStatementRepo) GetStatementsByUser(ctx context.Context, userId string) ([]*interop.Statement, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]*interop.Statement), args.Error(1)
}

func (m *MockStatementRepo) GetStatement(ctx context.Context, statement *interop.Statement) (*interop.Statement, error) {
	args := m.Called(ctx, statement)
	return args.Get(0).(*interop.Statement), args.Error(1)
}

func (m *MockStatementRepo) Save(statement *interop.Statement) (*interop.Statement, error) {
	args := m.Called(statement)
	return args.Get(0).(*interop.Statement), nil
}

func (m *MockStatementRepo) Delete(statement *interop.Statement) (*interop.Statement, error) {
	args := m.Called(statement)
	return args.Get(0).(*interop.Statement), args.Error(0)
}

func (m *MockStatementRepo) GetStatementDocContent(ctx context.Context, statement *interop.Statement) ([]byte, error) {
	args := m.Called(ctx, statement)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStatementRepo) FindOne(statement *interop.Statement) (*interop.Statement, error) {
	args := m.Called(statement)
	return args.Get(0).(*interop.Statement), args.Error(1)
}

func (m *MockStatementRepo) Query(statement *interop.Statement) []*interop.Statement {
	args := m.Called(statement)
	return args.Get(0).([]*interop.Statement)
}

func (m *MockStatementRepo) Get(id uint) {
	m.Called(id)
	return
}

var _ interfaces.IStatementsRepository = (*MockStatementRepo)(nil)
