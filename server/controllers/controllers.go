package controllers

import (
	"context"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type NotePurchaseAgreementController struct {
	interop.UnsafeNotePurchaseAgreementServiceServer
	notePurchaseAgreementService interfaces.IAgreementsService
}

func (c *NotePurchaseAgreementController) GetNotePurchaseAgreementDoc(ctx context.Context, req *interop.GetNotePurchaseAgreementDocRequest) (*interop.GetNotePurchaseAgreementDocResponse, error) {

	response := &interop.GetNotePurchaseAgreementDocResponse{}

	doc, err := c.notePurchaseAgreementService.GetNotePurchaseAgreementDocContent(context.TODO(), req.Payload)

	if err != nil {
		return response, err
	}

	response.FileContent = doc

	return response, nil
}

func NewNotePurchaseAgreementController(documentContentService interfaces.IAgreementsService) *NotePurchaseAgreementController {
	return &NotePurchaseAgreementController{
		notePurchaseAgreementService: documentContentService,
	}
}

var _ interop.NotePurchaseAgreementServiceServer = (*NotePurchaseAgreementController)(nil)
