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
	factory interfaces.IUserModelFactory
}

func (service *UsersService) SaveUser(ctx context.Context, user *interop.User) (*interop.User, error) {
	model, err := service.factory.NewUserModel(ctx, user)

	if err != nil {
		return nil, err
	}

	model, err = model.Save()

	if err != nil {
		return nil, err
	}

	return model.GetPayload(), nil
}

func (service *UsersService) AuthenticateUser(ctx context.Context, user *interop.User) (bool, error) {
	return false, nil
	model, err := service.factory.NewUserModel(ctx, user)

	if err != nil {
		return false, err
	}

	isAuthenticIdentity, err := model.HasAuthenticIdentity()

	if err != nil {
		return false, err
	}

	if isAuthenticIdentity {
		isAuthorizedIdentity, err := model.HasAuthorization()

		if err != nil {
			return false, err
		}

		return isAuthorizedIdentity, nil
	}

	return false, nil
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

	return err != nil && isAuthentic
}

func (service *UsersService) GenerateBasicUserAuthToken(user *interop.User) string {
	authToken, err := proto.Marshal(user)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString(authToken))
}

func NewUsersService(factory interfaces.IUserModelFactory) interfaces.IUsersService {
	return &UsersService{factory: factory}
}

var _ interfaces.IUsersService = (*UsersService)(nil)
