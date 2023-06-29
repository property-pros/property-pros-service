package controllers_test

import (
	"context"
	"testing"

	"github.com/vireocloud/property-pros-service/constants"
	"github.com/vireocloud/property-pros-service/interop"
	"google.golang.org/grpc"
)

func TestNotePurchaseAgreementController_EndToEnd(t *testing.T) {
    // Set up a connection to the gRPC server
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        t.Fatalf("Failed to connect to gRPC server: %v", err)
    }
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
        Payload: saveRes.GetPayload().GetId(),
    })
    if err != nil {
        t.Fatalf("Failed to get note purchase agreement: %v", err)
    }

    // Assert the saved and retrieved note purchase agreements match
    if getRes.GetPayload().GetId() != testNotePurchaseAgreement.GetId() {
        t.Errorf("Saved and retrieved note purchase agreements do not match.")
    }
}
