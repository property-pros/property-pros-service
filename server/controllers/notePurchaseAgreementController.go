package controllers

import (
	"context"
	"fmt"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

const (
	GRPC_REGISTRATION_METHOD string = "/api.note_purchase_agreement.v1.NotePurchaseAgreementService/SaveNotePurchaseAgreement"
)

type NotePurchaseAgreementController struct {
	interop.UnimplementedNotePurchaseAgreementServiceServer

	authService                  interfaces.IUsersService
	notePurchaseAgreementService interfaces.IAgreementsService
}

func notePurchaseAgreementToRecordResult(agreement *interop.NotePurchaseAgreement) *interop.RecordResultPayload {
	return &interop.RecordResultPayload{
		Id:        agreement.Id,
		CreatedOn: agreement.CreatedOn,
	}
}

func notePurchaseAgreementListToRecordCollection(result []interfaces.IAgreementModel) *interop.RecordColection {
	payload := []*interop.RecordResultPayload{}
	recordCollection := &interop.RecordColection{Payload: payload}

	for _, agreement := range result {
		payload = append(payload, notePurchaseAgreementToRecordResult(agreement.GetPayload()))
	}

	return recordCollection
}

func (c *NotePurchaseAgreementController) GetNotePurchaseAgreements(ctx context.Context, req *interop.GetNotePurchaseAgreementsRequest) (*interop.GetNotePurchaseAgreementsResponse, error) {

	response := &interop.GetNotePurchaseAgreementsResponse{}

	result, err := c.notePurchaseAgreementService.GetNotePurchaseAgreements(ctx)

	if err != nil {
		return response, err
	}

	response.Payload = notePurchaseAgreementListToRecordCollection(result)

	return response, nil
}

func (c *NotePurchaseAgreementController) GetNotePurchaseAgreement(ctx context.Context, req *interop.GetNotePurchaseAgreementRequest) (*interop.GetNotePurchaseAgreementResponse, error) {

	response := &interop.GetNotePurchaseAgreementResponse{}

	result, err := c.notePurchaseAgreementService.GetNotePurchaseAgreement(ctx, req.GetPayload())

	if err != nil {
		return response, err
	}

	response.Payload = result

	return response, nil
}

func (c *NotePurchaseAgreementController) SaveNotePurchaseAgreement(ctx context.Context, req *interop.SaveNotePurchaseAgreementRequest) (response *interop.SaveNotePurchaseAgreementResponse, errResult error) {

	response = &interop.SaveNotePurchaseAgreementResponse{}

	result, err := c.notePurchaseAgreementService.Save(ctx, req.Payload)
	fmt.Printf("user id: %v", result.User.Id)
	if err != nil {
		return response, err
	}

	response.Payload = notePurchaseAgreementToRecordResult(result)

	return response, nil
}

func (c *NotePurchaseAgreementController) GetNotePurchaseAgreementDoc(ctx context.Context, req *interop.GetNotePurchaseAgreementDocRequest) (response *interop.GetNotePurchaseAgreementDocResponse, errResult error) {

	response = &interop.GetNotePurchaseAgreementDocResponse{}

	doc, returnErr := c.notePurchaseAgreementService.GetNotePurchaseAgreementDocContent(ctx, req.GetPayload())

	if returnErr != nil {
		return response, returnErr
	}

	response.FileContent = doc

	return response, returnErr
}

func NewNotePurchaseAgreementController(notePurchaseAgreementService interfaces.IAgreementsService, authService interfaces.IUsersService) *NotePurchaseAgreementController {
	return &NotePurchaseAgreementController{
		notePurchaseAgreementService: notePurchaseAgreementService,
		authService:                  authService,
	}
}

var _ interop.NotePurchaseAgreementServiceServer = (*NotePurchaseAgreementController)(nil)
