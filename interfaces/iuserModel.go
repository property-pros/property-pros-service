package interfaces

import "github.com/vireocloud/property-pros-service/interop"

type IUserModel interface {
	IBaseModel[interop.User]
	Save() (IUserModel, error)
	MatchCredentials(IUserModel) (bool, error)
	HasAuthenticIdentity() (bool, error)
	HasAuthorization() (bool, error)
	GetPassword() string
}
