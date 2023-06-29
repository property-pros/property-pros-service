package controllers

import (
	"context"

	"github.com/vireocloud/property-pros-service/common"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatementController struct {
	interop.UnimplementedStatementServiceServer

	authService      interfaces.IUsersService
	statementService interfaces.IAgreementsService
	statementsRepo   interfaces.IStatementsRepository
	logger           *common.Logger
}

func statementToRecordResult(statement *interop.Statement) *interop.RecordResultPayload {
	return &interop.RecordResultPayload{
		Id: statement.Id,
		// CreatedOn: statement.CreatedOn,
	}
}

func statementListToRecordCollection(result []interfaces.IAgreementModel) *interop.RecordColection {
	payload := []*interop.RecordResultPayload{}
	recordCollection := &interop.RecordColection{Payload: payload}

	// for _, agreement := range result {
	// payload = append(payload, statementToRecordResult(agreement.GetPayload()))
	// }

	return recordCollection
}
func (c *StatementController) GetStatements(ctx context.Context, req *interop.GetStatementsRequest) (*interop.GetStatementsResponse, error) {
	c.logger.Info("Received GetStatements request")

	query := &interop.Statement{UserId: req.UserId}

	results, err := c.statementsRepo.Query(query)
	if err != nil {
		return nil, err
	}

	response := &interop.GetStatementsResponse{
		Payload: &interop.StatementsPayload{
			Statements: results,
		},
	}

	c.logger.Printf("Returning GetStatements response - StatementsCount: %v", len(results))

	return response, nil
}

func (c *StatementController) GetStatements(ctx context.Context, req *interop.GetStatementsRequest) (*interop.GetStatementsResponse, error) {
	c.logger.Info("Received GetStatements request")

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "UserId is required")
	}

	response := &interop.GetStatementsResponse{}

	query := &interop.Statement{UserId: req.GetUserId()}

	results := c.statementsRepo.Query(query)

	response.Payload.Statements = results

	c.logger.Info("Returning GetStatements response - StatementsCount: %v", len(results))

	return response, nil
}


// func (c *StatementController) GetStatement(ctx context.Context, req *interop.GetStatementRequest) (*interop.GetStatementResponse, error) {

// 	response := &interop.GetStatementResponse{}

// 	result, err := c.statementService.GetStatement(ctx, req.GetPayload())

// 	if err != nil {
// 		return response, err
// 	}

// 	response.Payload = result

// 	return response, nil
// }

// func (c *StatementController) SaveStatement(ctx context.Context, req *interop.SaveStatementRequest) (response *interop.SaveStatementResponse, errResult error) {

// 	response = &interop.SaveStatementResponse{}

// 	result, err := c.statementService.Save(ctx, req.Payload)

// 	if err != nil {
// 		return response, err
// 	}

// 	response.Payload = statementToRecordResult(result)

// 	return response, nil
// }

// func (c *StatementController) GetStatementDoc(ctx context.Context, req *interop.GetStatementDocRequest) (response *interop.GetStatementDocResponse, errResult error) {

// 	response = &interop.GetStatementDocResponse{}

// 	doc, returnErr := c.statementService.GetStatementDocContent(ctx, req.GetPayload())

// 	if returnErr != nil {
// 		return response, returnErr
// 	}

// 	response.FileContent = doc

// 	return response, returnErr
// }

func NewStatementController(statementRepo interfaces.IStatementsRepository) *StatementController {
	return &StatementController{
		statementsRepo: statementRepo,
	}
}

var _ interop.StatementServiceServer = (*StatementController)(nil)
