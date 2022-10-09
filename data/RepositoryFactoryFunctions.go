package data

import (
	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
	"gorm.io/gorm"
)

func NewRepository[T any, PT RepositoryModelConstraint[T]](db *gorm.DB) interfaces.IRepository[T] {
	return NewGormRepository[T, PT](db)
}

func NewUsersRepository(db *gorm.DB) interfaces.IUsersRepository {
	return NewRepository[interop.User](db)
}

func NewAgreementsRepository(db *gorm.DB) interfaces.IAgreementsRepository {
	return NewRepository[interop.NotePurchaseAgreement](db)
}
