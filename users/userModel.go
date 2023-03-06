package users

import (
	"context"
	"log"

	"github.com/vireocloud/property-pros-service/common"
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

type UserModel2 struct {
	gateway interfaces.IUsersGateway
}
type UserModel struct {
	*interop.User
	*common.BaseModel[interop.User]
	gateway interfaces.IUsersGateway
}

func (model *UserModel) GetContext() context.Context {
	return model.Context
}

func (model *UserModel) Save() (interfaces.IUserModel, error) {
	return nil, nil
	// return model.gateway.SaveUser(model)
}

func (model *UserModel) HasAuthenticIdentity() (bool, error) {
	userModel, err := model.gateway.GetUserByUsername(model)

	if err != nil {
		return false, nil
	}

	return model.MatchCredentials(userModel)
}

func (model *UserModel) MatchCredentials(identity interfaces.IUserModel) (bool, error) {
	return model.GetPassword() == identity.GetPassword(), nil
}

func (model *UserModel) HasAuthorization() (bool, error) {
	log.Default().Println("model: ", model)
	return model.GetId() != "", nil
}

func NewUserModel(gateway interfaces.IUsersGateway) *UserModel {
	user := &interop.User{}
	baseModel := &common.BaseModel[interop.User]{Payload: user, Context: context.Background()}
	return &UserModel{
		gateway:   gateway,
		User:      user,
		BaseModel: baseModel,
	}
}
