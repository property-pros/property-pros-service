package controllers

import (
	"context"
	"errors"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	GRPC_REGISTRATION_METHOD string = "/notepurchaseagreement.NotePurchaseAgreementService/SaveNotePurchaseAgreement"
	GRPC_AUTH_METHOD         string = "/auth.AuthenticationService/AuthenticateUser"
)

type PropertyProsApiController struct {
	interop.UnimplementedAuthenticationServiceServer
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

func (c *PropertyProsApiController) AuthenticateUser(ctx context.Context, req *interop.AuthenticateUserRequest) (*interop.AuthenticateUserResponse, error) {

	response := &interop.AuthenticateUserResponse{}

	isAuthentic, err := c.authService.AuthenticateUser(ctx, req.Payload)

	if err != nil {
		return response, err
	}

	response.Authenticated = isAuthentic

	// We want to extract metadata from the incomming context.
	// We dont create a new context since we dont wanna overwrite old metadata
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("could not grab metadata from context")
	}
	// Set authorization  metadata for the client to send in subsequent requests
	meta.Set("authorization", c.authService.GenerateBasicUserAuthToken(req.Payload))
	// Metadata is sent on its own, so we need to send the header. There is also something called Trailer
	grpc.SendHeader(ctx, meta)

	return response, nil
}

func NewNotePurchaseAgreementController(notePurchaseAgreementService interfaces.IAgreementsService, authService interfaces.IUsersService) *PropertyProsApiController {
	return &PropertyProsApiController{
		notePurchaseAgreementService: notePurchaseAgreementService,
		authService:                  authService,
	}
}

var _ interop.NotePurchaseAgreementServiceServer = (*PropertyProsApiController)(nil)

var _ interop.AuthenticationServiceServer = (*PropertyProsApiController)(nil)
