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
	GRPC_AUTH_METHOD string = "/auth.AuthenticationService/AuthenticateUser"
)

type AuthController struct {
	interop.UnimplementedAuthenticationServiceServer

	authService                  interfaces.IUsersService
	notePurchaseAgreementService interfaces.IAgreementsService
}

func (c *AuthController) AuthenticateUser(ctx context.Context, req *interop.AuthenticateUserRequest) (*interop.AuthenticateUserResponse, error) {

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

func NewAuthController(notePurchaseAgreementService interfaces.IAgreementsService, authService interfaces.IUsersService) *AuthController {
	return &AuthController{
		notePurchaseAgreementService: notePurchaseAgreementService,
		authService:                  authService,
	}
}

var _ interop.AuthenticationServiceServer = (*AuthController)(nil)
