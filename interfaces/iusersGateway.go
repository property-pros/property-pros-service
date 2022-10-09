package interfaces

type IUsersGateway interface {
	SaveUser(user IUserModel) (IUserModel, error)
	GetUserByUsername(user IUserModel) (IUserModel, error)
}
