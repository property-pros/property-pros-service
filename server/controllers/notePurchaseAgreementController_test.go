// +build integration,test.paniconexit0

package controllers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vireocloud/property-pros-service/config"
	"github.com/vireocloud/property-pros-service/constants"
	"github.com/vireocloud/property-pros-service/interop"
	"google.golang.org/grpc"
)

func TestNotePurchaseAgreementController_EndToEnd(t *testing.T) {
    
	config, err := config.NewConfig() 
	// with loading config
	require.NoError(t, err)
	// Create a connection to the gRPC server
	conn, err := grpc.Dial(config.DocumentContentProviderSource, grpc.WithInsecure())
	require.NoError(t, err)

    defer conn.Close()

    // Create a client for the NotePurchaseAgreementService
    client := interop.NewNotePurchaseAgreementServiceClient(conn)

    // Create test data
    testNotePurchaseAgreement := &interop.NotePurchaseAgreement{
        Id:        "1",
        CreatedOn: "2020-01-01T00:00:00Z",
    }

    // Save a new note purchase agreement
    ctx := context.WithValue(context.Background(), constants.UserIdKey, "1")
    saveRes, err := client.SaveNotePurchaseAgreement(ctx, &interop.SaveNotePurchaseAgreementRequest{
        Payload: testNotePurchaseAgreement,
    })
    if err != nil {
        t.Fatalf("Failed to save note purchase agreement: %v", err)
    }
	
       // Get the newly saved note purchase agreement
    getRes, err := client.GetNotePurchaseAgreement(ctx, &interop.GetNotePurchaseAgreementRequest{
        Payload: &interop.RecordRequestPayload{
            Id: saveRes.GetPayload().GetId(),
        },
    }) 


    if err != nil {
        t.Fatalf("Failed to get note purchase agreement: %v", err)
    }

    // Assert the saved and retrieved note purchase agreements match
    if getRes.GetPayload().GetId() != testNotePurchaseAgreement.GetId() {
        t.Errorf("Saved and retrieved note purchase agreements do not match.")
    }
}

// Test GetNotePurchaseAgreementDoc
func TestNotePurchaseAgreementController_GetNotePurchaseAgreementDoc(t *testing.T) {
    // Set up a connection to the gRPC server
	config, err := config.NewConfig() 
	// with loading config
	require.NoError(t, err)
	// Create a connection to the gRPC server
	conn, err := grpc.Dial(config.DocumentContentProviderSource, grpc.WithInsecure())
	require.NoError(t, err)

    defer conn.Close()

    // Create a client for the NotePurchaseAgreementService
    client := interop.NewNotePurchaseAgreementServiceClient(conn)

    // Create test data
    testNotePurchaseAgreementDoc := &interop.NotePurchaseAgreement{
        Id:        "1",
        CreatedOn: "2020-01-01T00:00:00Z",
        FileContent: []byte("test file content"),
    }

       // Get the newly saved note purchase agreement doc
    getRes, err := client.GetNotePurchaseAgreementDoc(context.TODO(), &interop.GetNotePurchaseAgreementDocRequest{
        Payload: testNotePurchaseAgreementDoc,
    })
    if err != nil {
        t.Fatalf("Failed to get note purchase agreement doc: %v", err)
    }

    // Assert the saved and retrieved note purchase agreement docs match
    if !bytes.Equal(getRes.FileContent, testNotePurchaseAgreementDoc.FileContent) {
        t.Errorf("Saved and retrieved note purchase agreement docs do not match.")
    }
}

// Test GetNotePurchaseAgreements
func TestNotePurchaseAgreementController_GetNotePurchaseAgreements(t *testing.T) {
    // Set up a connection to the gRPC server
	config, err := config.NewConfig() 
	// with loading config
	require.NoError(t, err)
	// Create a connection to the gRPC server
	conn, err := grpc.Dial(config.DocumentContentProviderSource, grpc.WithInsecure())
	require.NoError(t, err)

    defer conn.Close()

    // Create a client for the NotePurchaseAgreementService
    client := interop.NewNotePurchaseAgreementServiceClient(conn)

    // Create test data
    testNotePurchaseAgreement := &interop.NotePurchaseAgreement{
        Id:        "1",
        CreatedOn: "2020-01-01T00:00:00Z",
    }

    ctx := context.WithValue(context.Background(), constants.UserIdKey, "1")
    // Save a new note purchase agreement
    saveRes, err := client.SaveNotePurchaseAgreement(ctx, &interop.SaveNotePurchaseAgreementRequest{
        Payload: testNotePurchaseAgreement,
    })
    if err != nil {
        t.Fatalf("Failed to save note purchase agreement: %v", err)
    }

    // Get the newly saved note purchase agreements
    getRes, err := client.GetNotePurchaseAgreements(ctx, &interop.GetNotePurchaseAgreementsRequest{})
    if err != nil {
        t.Fatalf("Failed to get note purchase agreements: %v", err)
    }

    // Assert the saved note purchase agreement is in the retrieved list
    var found bool
    for _, npa := range getRes.GetPayload().GetPayload() {
        if npa.GetId() == saveRes.GetPayload().GetId() {
            found = true
            break
        }
    }
    if !found {
        t.Errorf("Saved note purchase agreement not found in retrieved list.")
    }
}




