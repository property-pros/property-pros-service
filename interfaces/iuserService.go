package interfaces

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
)

type IUsersService interface {
	SaveUser(context.Context, *interop.User) (*interop.User, error)
	AuthenticateUser(context.Context, *interop.User) (bool, error)
	IsValidToken(ctx context.Context, token string) bool
	GenerateBasicUserAuthToken(*interop.User) string
}
