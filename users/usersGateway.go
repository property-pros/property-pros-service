package users

import (
	"fmt"

	"github.com/vireocloud/property-pros-service/interfaces"
)

type UsersGateway struct {
	repo    interfaces.IUsersRepository
	factory interfaces.IUserModelFactory
}

func (gateway *UsersGateway) GetUserByUsername(user interfaces.IUserModel) (interfaces.IUserModel, error) {
	return nil, fmt.Errorf("GetUserByUsername Unimplemented")
}

func (gateway *UsersGateway) SaveUser(user interfaces.IUserModel) (interfaces.IUserModel, error) {

	result, err := gateway.repo.Save(user.GetPayload())

	if err != nil {
		return nil, err
	}

	return gateway.factory.NewUserModel(user.GetContext(), result)
}

func (gateway *UsersGateway) CreateNewUser(user interfaces.IUserModel) (interfaces.IUserModel, error) {

	result, err := gateway.repo.Save(user.GetPayload())

	if err != nil {
		return nil, err
	}

	return gateway.factory.NewUserModel(user.GetContext(), result)
}

func (gateway *UsersGateway) UpdateUser(user interfaces.IUserModel) (interfaces.IUserModel, error) {
	return nil, nil
}

func NewUsersGateway(repo interfaces.IUsersRepository, factory interfaces.IUserModelFactory) interfaces.IUsersGateway {
	return &UsersGateway{
		repo:    repo,
		factory: factory,
	}
}

var _ interfaces.IUsersGateway = (*UsersGateway)(nil)
