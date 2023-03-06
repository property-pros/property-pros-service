package users

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vireocloud/property-pros-service/data"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
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

func (gateway *UsersGateway) SaveUser(ctx context.Context, agreement *interop.NotePurchaseAgreement) (*interop.User, error) {
	userData := data.User{
		Id:             uuid.New().String(),
		FirstName:      agreement.GetFirstName(),
		LastName:       agreement.GetLastName(),
		DateOfBirth:    agreement.GetDateOfBirth(),
		EmailAddress:   agreement.User.GetEmailAddress(),
		Password:       agreement.User.GetPassword(),
		HomeAddress:    agreement.GetHomeAddress(),
		PhoneNumber:    agreement.GetPhoneNumber(),
		SocialSecurity: agreement.GetSocialSecurity(),
	}

	user, err := gateway.repo.Save(&userData)

	if err != nil {
		return nil, err
	}

	return &interop.User{
		Id:           user.Id,
		EmailAddress: user.EmailAddress,
		Password:     user.Password,
	}, nil
}

func (gateway *UsersGateway) CreateNewUser(user data.User) (*data.User, error) {
	user.Id = uuid.New().String()
	return nil, nil
	// return gateway.SaveUser(user)
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
