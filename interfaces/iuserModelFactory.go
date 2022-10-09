package interfaces

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
)

type IUserModelFactory interface {
	NewUserModel(context.Context, *interop.User) (IUserModel, error)
}
