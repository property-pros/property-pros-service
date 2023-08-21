package controllers

import (
	"context"
	"fmt"

	"github.com/vireocloud/property-pros-service/common"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatementController struct {
	interop.UnimplementedStatementServiceServer

	statementsRepo  interfaces.IStatementsRepository
	logger          *common.Logger
	documentService interfaces.IDocumentContentService
}

func NewStatementController(statementRepo interfaces.IStatementsRepository, documentService interfaces.IDocumentContentService) *StatementController {
	return &StatementController{
		statementsRepo:  statementRepo,
		documentService: documentService,
	}
}

func (c *StatementController) GetStatements(ctx context.Context, req *interop.GetStatementsRequest) (*interop.GetStatementsResponse, error) {
	c.logger.Info(fmt.Sprintf("Received GetStatements request: %v", req.String()))

	userId, err := GetUserIdFromContext(ctx)
	
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	query := &interop.Statement{UserId: userId}

	results := c.statementsRepo.Query(query)

	for _, statement := range results {

		doc, err := c.documentService.BuildStatement(ctx, statement)

		if err != nil {
			fmt.Printf("Error building statement doc: %v", err)
			return nil, err
		}

		statement.Document = doc.GetDocContent()
	}

	response := &interop.GetStatementsResponse{
		Payload: &interop.StatementsPayload{
			Statements: results,
		},
	}

	c.logger.Info("Returning GetStatements response - StatementsCount: %v", len(results))

	return response, nil
}

var _ interop.StatementServiceServer = (*StatementController)(nil)
