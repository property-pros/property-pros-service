package users

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"google.golang.org/protobuf/proto"
)

type UsersService struct {
	userGateway interfaces.IUsersGateway
}

func (service *UsersService) AuthenticateUser(ctx context.Context, user *interop.User) (string, error) {

	usr, err := service.userGateway.GetUserByUsername(user.EmailAddress)
	if err != nil {
		return "", err
	}

	if user.Password != usr.User.Password {
		return "", nil
	}

	return usr.User.Id, nil
}

func (service *UsersService) UserIdIfValidToken(ctx context.Context, token string) string {
	payload := &interop.User{}

	authToken, err := base64.StdEncoding.DecodeString(strings.Replace(token, "Basic ", "", 1))

	if err != nil {
		return ""
	}

	err = proto.Unmarshal(authToken, payload)

	if err != nil {
		return ""
	}

	usrID, err := service.AuthenticateUser(ctx, payload)
	if err != nil || usrID == "" {
		return ""
	}

	return usrID
}

func (service *UsersService) GenerateBasicUserAuthToken(user *interop.User) string {
	authToken, err := proto.Marshal(user)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString(authToken))
}

func NewUsersService(userGateway interfaces.IUsersGateway) interfaces.IUsersService {
	fmt.Printf("NewUsersService - user gateway: %v\r\n", userGateway)
	return &UsersService{userGateway: userGateway}
}

var _ interfaces.IUsersService = (*UsersService)(nil)
