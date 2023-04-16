package interfaces

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
)

type IUsersService interface {
	AuthenticateUser(context.Context, *interop.User) (string, error)
	UserIdIfValidToken(ctx context.Context, token string) string
	GenerateBasicUserAuthToken(*interop.User) string
}
