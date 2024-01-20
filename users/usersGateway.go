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

func (gateway *UsersGateway) GetUserByUsername(userEmail string) (*interop.NotePurchaseAgreement, error) {
	users := gateway.repo.Query(&data.User{
		EmailAddress: userEmail,
	})
	if len(users) == 0 {
		return nil, errors.New("no user found")
	}

	return &interop.NotePurchaseAgreement{
		FirstName:   users[0].FirstName,
		LastName:    users[0].LastName,
		DateOfBirth: users[0].DateOfBirth,
		User: &interop.User{
			EmailAddress: users[0].EmailAddress,
			Password:     users[0].Password,
			Id:           users[0].Id,
		},
		HomeAddress:    users[0].HomeAddress,
		PhoneNumber:    users[0].PhoneNumber,
		SocialSecurity: users[0].SocialSecurity,
	}, nil
}

func (gateway *UsersGateway) GetUser(user data.User) (*data.User, error) {
	return gateway.repo.FindOne(&user)
}

func (gateway *UsersGateway) SaveUser(ctx context.Context, agreement *interop.NotePurchaseAgreement) (*interop.NotePurchaseAgreement, error) {

	userId := agreement.GetUser().GetId()

	if userId == "" {
		userId = uuid.New().String()
	}

	userData := data.User{
		Id:             userId,
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

	agreementResultValue := *agreement

	agreementResult := &agreementResultValue

	agreementResult.User.Id = user.Id

	return agreement, nil
}

func (gateway *UsersGateway) CreateNewUser(user data.User) (*data.User, error) {
	user.Id = uuid.New().String()
	return nil, nil
	// return gateway.SaveUser(user)
}

func (gateway *UsersGateway) UpdateUser(user data.User) (data.User, error) {
	return data.User{}, nil
}

func NewUsersGateway(repo interfaces.IRepository[data.User]) interfaces.IUsersGateway {
	return &UsersGateway{
		repo: repo,
	}
}
