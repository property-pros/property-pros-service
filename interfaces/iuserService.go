package interfaces

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
)

type IUsersService interface {
	AuthenticateUser(context.Context, *interop.User) (bool, error)
	IsValidToken(ctx context.Context, token string) bool
	GenerateBasicUserAuthToken(*interop.User) string
}
