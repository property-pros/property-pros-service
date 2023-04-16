package data

import (
	"gorm.io/gorm"
	//according to gorm docs, this import is necessary
	"errors"
	"log"

	_ "gorm.io/driver/postgres"

	i "github.com/vireocloud/property-pros-service/interfaces"
)

type RepositoryModelConstraint[T any] interface {
	i.IModelPayload
	*T
}

type GormRepository[T any, PT RepositoryModelConstraint[T]] struct {
	i.IRepository[T]
	db *gorm.DB
}

var ID_FIELD_NAME = "Id"

func (repo *GormRepository[T, PT]) SetDb(db *gorm.DB) {
	repo.db = db
}

func (repo *GormRepository[T, PT]) Save(payload *T) (*T, error) {
	var model *T = payload

	modelResult := repo.db.Debug().Model(payload)
	err := modelResult.Save(model).Error
	return model, err
}

func (repo *GormRepository[T, PT]) Create(payload *T, query *T) (*T, error) {

	var model *T = payload

	modelResult := repo.db.Model(payload)

	err := modelResult.Create(model).Error

	return model, err
}

func (repo *GormRepository[T, PT]) Update(payload *T, query *T) (*T, error) {
	model := repo.db.Debug().Model(payload)
	modelTransaction := model.Begin()
	// modelTransaction.
	modelUpdates := modelTransaction.Set(PT(payload).GetId(), payload)
	// modelUpdates := modelQuery.Updates(payload)
	modelCommitted := modelUpdates.Commit()
	err := modelCommitted.Scan(payload).Error

	return payload, err
}

func (repo *GormRepository[T, PT]) FindOne(payload *T) (*T, error) {
	results := []*T{}

	err := repo.db.Debug().Model(payload).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	resultCount := len(results)

	if resultCount > 1 {
		return nil, errors.New("more than one result returned for query in FindOne()")
	}

	if resultCount == 0 {
		return nil, errors.New("zero results returned for query in FindOne()")
	}

	return results[0], nil
}

func (repo *GormRepository[T, PT]) Query(query *T) []*T {

	var results interface{} = []*T{}

	// whereResult := repo.db.Where()
	repo.db.Debug().Find(&results, query)

	return results.([]*T)
}

func (repo *GormRepository[T, PT]) Delete(query *T) (*T, error) {
	log.Printf("delete entity: %v", query)

	repo.db.Debug().Model(query).Delete(query)

	return query, nil
}

func NewGormRepository[T any, PT RepositoryModelConstraint[T]](db *gorm.DB) i.IRepository[T] {
	return &GormRepository[T, PT]{
		db: db,
	}
}
