package users

import (
	"fmt"

	"github.com/vireocloud/property-pros-service/data"
	"github.com/vireocloud/property-pros-service/interfaces"
)

type UsersGateway struct {
	repo    interfaces.IRepository[data.User]
	factory interfaces.IUserModelFactory
}

func (gateway *UsersGateway) GetUserByUsername(user interfaces.IUserModel) (interfaces.IUserModel, error) {
	return nil, fmt.Errorf("GetUserByUsername Unimplemented")
}

func (gateway *UsersGateway) SaveUser(user data.User) (data.User, error) {
	fmt.Println("here in gateway")
	result, err := gateway.repo.Save(&user)

	if err != nil {
		return data.User{}, err
	}

	return *result, nil
}

func (gateway *UsersGateway) CreateNewUser(user data.User) (data.User, error) {
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
