package interfaces

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
)

type IUsersGateway interface {
	SaveUser(ctx context.Context, user *interop.NotePurchaseAgreement) (*interop.User, error)
	GetUserByUsername(user IUserModel) (IUserModel, error)
}
