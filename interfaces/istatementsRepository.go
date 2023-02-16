package interfaces

import "github.com/vireocloud/property-pros-service/interop"

type IStatementsRepository interface {
	IRepository[interop.Statement]
}
