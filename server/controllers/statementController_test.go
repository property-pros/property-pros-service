//go:build integration
// +build integration

package controllers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	statement "github.com/vireocloud/property-pros-sdk/api/statement/v1"
	"google.golang.org/grpc"
)

func TestGetStatements(t *testing.T) {
	// Create a connection to the gRPC server
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	require.NoError(t, err)

	// Create a client from the connection
	client := statement.NewStatementServiceClient(conn)

	// Create the request
	req := &statement.GetStatementsRequest{
		UserId: "user1",
	}

	// Call the RPC and get the response
	response, err := client.GetStatements(context.Background(), req)
	require.NoError(t, err)

	// Assert on the response
	assert.Equal(t, 2, len(response.GetPayload().GetStatements()))
	assert.Equal(t, "statement1", response.GetPayload().Statements[0].GetId())
}

func TestGetStatementDoc(t *testing.T) {
	// Create a connection to the gRPC server
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	require.NoError(t, err)

	// Create a client from the connection
	client := statement.NewStatementServiceClient(conn)

	// Create the request
	req := &statement.GetStatementDocRequest{
		Payload: &statement.Statement{
			Id: "statement1",
		},
	}

	// Call the RPC and get the response
	response, err := client.GetStatementDoc(context.Background(), req)
	require.NoError(t, err)

	// Assert on the response
	assert.Equal(t, "statement1.pdf", response.GetDocument())
}
