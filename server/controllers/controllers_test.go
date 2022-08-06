package controllers

import (
	"testing"

	propertyProsApi "github.com/vireocloud/property-pros-docs/generated/notePurchaseAgreement"
)

func TestGetNotePurchaseAgreementDoc(t *testing.T) {
	controller := &NotePurchaseAgreementController{}

	result, err := controller.GetNotePurchaseAgreementDoc(nil, &propertyProsApi.GetNotePurchaseAgreementDocRequest{
		Payload: &propertyProsApi.NotePuchaseAgreement{
			FirstName: "John",
			LastName:  "smith",
		},
	})

	if err != nil {
		t.Error(err)
	}

	if result == nil {
		t.Errorf("Expected controller.GetNotePurchaseAgreementDoc to not return nil")
	}
}
