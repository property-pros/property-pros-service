package interfaces

import "github.com/vireocloud/property-pros-service/interop"

type IUsersRepository interface {
	IRepository[interop.User]
}
