package data

import (
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"gorm.io/gorm"
)

func NewRepository[T any, PT RepositoryModelConstraint[T]](db *gorm.DB) interfaces.IRepository[T] {
	return NewGormRepository[T, PT](db)
}

func NewUsersRepository(db *gorm.DB) interfaces.IRepository[User] {
	return NewRepository[User](db)
}

func NewAgreementsRepository(db *gorm.DB) interfaces.IRepository[NotePurchaseAgreement] {
	return NewRepository[NotePurchaseAgreement](db)
}

// func NewDocumentsRepository(db *gorm.DB) interfaces.IAgreementsRepository {
// 	return NewRepository[interop.Document](db)
// }
func NewStatementsRepository(db *gorm.DB) interfaces.IStatementsRepository {
	return NewRepository[interop.Statement](db)
}
