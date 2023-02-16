package controllers

import (
	"context"

	"github.com/vireocloud/property-pros-sdk/api/statement/v1"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type StatementController struct {
	interop.UnimplementedStatementServiceServer

	authService      interfaces.IUsersService
	statementService interfaces.IAgreementsService
	statementsRepo   interfaces.IStatementsRepository
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

	response := &interop.GetStatementsResponse{}

	query := &interop.Statement{UserId: req.GetUserId()}

	// response.Statements = c.statementsRepo.Query(query)

	response.Statements = []*statement.Statement{
		{UserId: query.UserId },
	}

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
