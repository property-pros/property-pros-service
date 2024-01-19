package controllers

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vireocloud/property-pros-sdk/api/auth/v1"
	"github.com/vireocloud/property-pros-sdk/api/common/v1"
	"github.com/vireocloud/property-pros-service/config"
	"github.com/vireocloud/property-pros-service/interop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type StatementControllerApiTestSuite struct {
	suite.Suite
	conn                      *grpc.ClientConn
	statementClient           interop.StatementServiceClient
	authClient                auth.AuthenticationServiceClient
	Config                    *config.Config
	authToken                 string
	testNotePurchaseAgreement *interop.NotePurchaseAgreement
	npaClient                 interop.NotePurchaseAgreementServiceClient
	testServiceEndpoint       string
}

func (s *StatementControllerApiTestSuite) SetupSuite() {
	var err error

	s.Config, err = config.NewConfig()

	if err != nil {
		s.T().Fatalf("Failed to load config: %v", err)
	}
	s.testServiceEndpoint = os.Getenv("TEST_SERVICE")

	s.conn, err = grpc.Dial(s.testServiceEndpoint, grpc.WithInsecure())
	if err != nil {
		s.T().Fatalf("Failed to connect: %v", err)
	}
	s.statementClient = interop.NewStatementServiceClient(s.conn)

	// Add NotePurchaseAgreementService client
	npaConn, err := grpc.Dial(s.testServiceEndpoint, grpc.WithInsecure())
	if err != nil {
		s.T().Fatalf("Failed to connect: %v", err)
	}
	s.npaClient = interop.NewNotePurchaseAgreementServiceClient(npaConn)

	// Create a new user
	ctx := context.Background()
	user := &interop.User{
		Id:           "user1",
		EmailAddress: "user1@example.com",
		Password:     "password",
	}
	req := &interop.SaveNotePurchaseAgreementRequest{
		Payload: &interop.NotePurchaseAgreement{
			User: user,
		},
	}
	response, err := s.npaClient.SaveNotePurchaseAgreement(ctx, req)
	if err != nil {
		s.T().Fatalf("SaveNotePurchaseAgreement failed: %v", err)
	}

	testNotePurchaseAgreementId := response.Payload.Id

	// Authenticate user to get auth token
	s.authClient = auth.NewAuthenticationServiceClient(s.conn)
	authReq := &auth.AuthenticateUserRequest{
		Payload: user,
	}

	headerMD := metadata.MD{}

	authResp, err := s.authClient.AuthenticateUser(ctx, authReq, grpc.Header(&headerMD))
	if err != nil {
		s.T().Fatalf("failed to call: %v", err)
	}

	s.authToken = headerMD["authorization"][0]

	s.T().Logf("authToken: %s", s.authToken)
	if err != nil {
		s.T().Fatalf("AuthenticateUser failed: %v", err)
	}

	if !authResp.IsAuthenticated {
		s.T().Fatalf("User not authenticated")
	}

	md := metadata.Pairs(
		"authorization", s.authToken,
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	notePurchaseAgreementResponse, err := s.npaClient.GetNotePurchaseAgreement(ctx, &interop.GetNotePurchaseAgreementRequest{
		Payload: &common.RecordRequestPayload{
			Id: testNotePurchaseAgreementId,
		},
	})

	if err != nil {
		s.T().Fatalf("GetNotePurchaseAgreement failed: %v", err)
	}

	s.testNotePurchaseAgreement = notePurchaseAgreementResponse.Payload
	s.T().Logf("testNotePurchaseAgreement: %+v", s.testNotePurchaseAgreement)
	s.T().Logf("testNotePurchaseAgreement user: %+v", s.testNotePurchaseAgreement.User)
}

func (s *StatementControllerApiTestSuite) TearDownSuite() {
	s.conn.Close()
}

func (s *StatementControllerApiTestSuite) TestGetStatements() {
	ctx := context.Background()
	md := metadata.Pairs(
		"authorization", s.authToken,
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	userId := os.Getenv("TEST_ACCOUNT_ID")

	req := &interop.GetStatementsRequest{
		UserId: userId,
	}

	resp, err := s.statementClient.GetStatements(ctx, req, grpc.Header(&metadata.MD{
		"authorization": []string{s.authToken},
	}))

	if err != nil {
		s.T().Fatalf("GetStatements failed: %v", err)
	}

	if len(resp.Payload.Statements) != 3 {
		s.T().Fatalf("Expected 3 statements, got %d", len(resp.Payload.Statements))
	}

	stmt := resp.Payload.Statements[0]
	if stmt.UserId != userId {
		s.T().Fatalf("Expected statement for %v, got %s", userId, stmt.UserId)
	}

	for _, stmt := range resp.Payload.Statements {
		if len(stmt.Document) == 0 {
			s.T().Errorf("Statement %v has blank document", stmt.Id)
		}
	}

}

func TestStatementControllerApiTestSuite(t *testing.T) {
	suite.Run(t, new(StatementControllerApiTestSuite))
}
