package controllers

import (
	"context"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

const (
	GRPC_REGISTRATION_METHOD string = "/notepurchaseagreement.NotePurchaseAgreementService/SaveNotePurchaseAgreement"
)

type NotePurchaseAgreementController struct {
	interop.UnimplementedNotePurchaseAgreementServiceServer

	authService                  interfaces.IUsersService
	notePurchaseAgreementService interfaces.IAgreementsService
}

func (c *PropertyProsApiController) SaveNotePurchaseAgreement(ctx context.Context, req *interop.SaveNotePurchaseAgreementRequest) (response *interop.SaveNotePurchaseAgreementResponse, errResult error) {

	response = &interop.SaveNotePurchaseAgreementResponse{}

	result, err := c.notePurchaseAgreementService.Save(ctx, req.Payload)

	if err != nil {
		return response, err
	}

	response.Id = result.GetId()

	return response, nil
}

func (c *PropertyProsApiController) GetNotePurchaseAgreementDoc(ctx context.Context, req *interop.GetNotePurchaseAgreementDocRequest) (response *interop.GetNotePurchaseAgreementDocResponse, errResult error) {

	response = &interop.GetNotePurchaseAgreementDocResponse{}

	doc, returnErr := c.notePurchaseAgreementService.GetNotePurchaseAgreementDocContent(ctx, req.Payload)

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

var _ interop.NotePurchaseAgreementServiceServer = (*PropertyProsApiController)(nil)
