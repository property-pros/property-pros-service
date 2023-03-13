package interfaces

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
)

type IUsersGateway interface {
	SaveUser(ctx context.Context, user *interop.NotePurchaseAgreement) (*interop.NotePurchaseAgreement, error)
	GetUserByUsername(userEmail string) (*interop.NotePurchaseAgreement, error)
}
