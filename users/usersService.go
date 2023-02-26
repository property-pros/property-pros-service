package users

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/vireocloud/property-pros-service/data"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"google.golang.org/protobuf/proto"
)

type UsersService struct {
	userGateway *UsersGateway
}

func (service *UsersService) SaveUser(ctx context.Context, user *interop.User) (*interop.User, error) {
	// model, err := service.factory.NewUserModel(ctx, user)

	// if err != nil {
	// 	fmt.Printf("%v error in SaveUser userService", err)
	// 	return nil, err
	// }
	// fmt.Println("saving user")
	// fmt.Printf("%v\n", model)
	// model, err = model.Save()

	// if err != nil {
	// 	return nil, err
	// }

	return &interop.User{}, nil
}

func (service *UsersService) AuthenticateUser(ctx context.Context, user *interop.User) (bool, error) {
	usr, err := service.userGateway.GetUserByUsername(data.User{
		EmailAddress: user.EmailAddress,
	})
	if err != nil {
		return false, err
	}

	if user.Password != usr.Password {
		return false, nil
	}

	return true, nil
}

func (service *UsersService) IsValidToken(ctx context.Context, token string) bool {
	payload := &interop.User{}

	authToken, err := base64.StdEncoding.DecodeString(strings.Replace(token, "Basic ", "", 1))

	if err != nil {
		return false
	}

	err = proto.Unmarshal(authToken, payload)

	if err != nil {
		return false
	}

	isAuthentic, err := service.AuthenticateUser(ctx, payload)

	return err == nil && isAuthentic
}

func (service *UsersService) GenerateBasicUserAuthToken(user *interop.User) string {
	authToken, err := proto.Marshal(user)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString(authToken))
}

func NewUsersService(userGateway *UsersGateway) interfaces.IUsersService {
	return &UsersService{userGateway: userGateway}
}

var _ interfaces.IUsersService = (*UsersService)(nil)
