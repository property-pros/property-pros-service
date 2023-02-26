package users

import (
	"errors"

	"github.com/vireocloud/property-pros-service/data"
	"github.com/vireocloud/property-pros-service/interfaces"
)

type UsersGateway struct {
	repo    interfaces.IRepository[data.User]
	factory interfaces.IUserModelFactory
}

func (gateway *UsersGateway) GetUserByUsername(user data.User) (*data.User, error) {
	users := gateway.repo.Query(&user)
	if len(users) == 0 {
		return nil, errors.New("no user found")
	}

	return users[0], nil
}

func (gateway *UsersGateway) GetUser(user data.User) (*data.User, error) {
	return gateway.repo.FindOne(&user)
}

func (gateway *UsersGateway) SaveUser(user data.User) (*data.User, error) {
	result, err := gateway.repo.Save(&user)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (gateway *UsersGateway) CreateNewUser(user data.User) (*data.User, error) {
	return gateway.SaveUser(user)
}

func (gateway *UsersGateway) UpdateUser(user data.User) (data.User, error) {
	return data.User{}, nil
}

func NewUsersGateway(repo interfaces.IRepository[data.User], factory interfaces.IUserModelFactory) *UsersGateway {
	return &UsersGateway{
		repo:    repo,
		factory: factory,
	}
}
