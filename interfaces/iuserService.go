package interfaces

import "context"

type IUserService interface {
	SaveUser(context.Context, IUserModel) (IUserModel, error)
}
