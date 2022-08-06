package data

import (
	"github.com/jinzhu/gorm"
	//according to gorm docs, this import is necessary
	"errors"
	"log"
	"reflect"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	i "github.com/vireocloud/property-pros-service/interfaces"
)

type GormRepository struct {
	i.IRepository
	db *gorm.DB
}

var ID_FIELD_NAME = "ID"

func (repo *GormRepository) SetDb(db *gorm.DB) {
	repo.db = db
}

func (repo *GormRepository) Save(payload interface{}, query interface{}) (interface{}, error) {
	results, err := repo.FindOne(query)

	if err != nil && results == nil {
		return repo.db.Model(query).Create(payload).Value, nil
	}

	return repo.db.Model(query).Updates(payload).Value, nil
}

func (repo *GormRepository) FindOne(payload interface{}) (interface{}, error) {

	emptyResultSlicePointer := repo.createEmptyResultSet(payload)

	repo.db.Where(ID_FIELD_NAME+" = ?", repo.getIDValue(payload)).Find(emptyResultSlicePointer)

	results := InterfaceSlice(emptyResultSlicePointer)

	resultCount := len(results)

	if resultCount > 1 {
		return results, errors.New("More than one result returned for query in FindOne()")
	}

	if resultCount == 0 {
		return nil, errors.New("Zero results returned for query in FindOne()")
	}

	return results[0], nil
}

func (repo *GormRepository) Query(query interface{}) []interface{} {

	emptyResultSlicePointer := repo.createEmptyResultSet(query)

	repo.db.Where(query).Find(emptyResultSlicePointer)

	return InterfaceSlice(emptyResultSlicePointer)
}

func (repo *GormRepository) Delete(query interface{}) (interface{}, error) {
	log.Printf("delete entity: %v", query)
	repo.db.Delete(query)

	return query, nil
}

func (repo *GormRepository) ensureValueObject(entity interface{}) interface{} {
	payloadType := reflect.TypeOf(entity)

	if payloadType.Kind() == reflect.Ptr {
		return reflect.ValueOf(entity).Elem().Interface()
	}

	return entity
}

func (repo *GormRepository) Initialize() error {

	db, err := gorm.Open("postgres", "host=myhost port=myport user=gorm dbname=gorm password=mypassword")

	defer db.Close()

	repo.db = db

	return err
}

func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice).Elem()
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func (repo *GormRepository) GetQueryAndPayload(payload interface{}) (newPayload interface{}, query interface{}) {

	var payloadValueInfo reflect.Value
	var newQueryObject reflect.Value

	payloadType := reflect.TypeOf(payload)

	payloadIsPtr := payloadType.Kind() == reflect.Ptr

	if payloadIsPtr {
		payloadType = payloadType.Elem()

		payloadValueInfo = reflect.ValueOf(payload).Elem()
	} else {
		payloadValueInfo = reflect.ValueOf(payload)
	}

	newQueryObject = reflect.New(payloadType).Elem()

	idField := payloadValueInfo.FieldByName(ID_FIELD_NAME)

	newQueryObject.FieldByName(ID_FIELD_NAME).SetUint(idField.Uint())

	if payloadIsPtr {
		idField.SetUint(0)
	}

	if payloadIsPtr {
		payloadValueInfo = payloadValueInfo.Addr()
		newQueryObject = newQueryObject.Addr()
	}

	return payloadValueInfo.Interface(), newQueryObject.Interface()
}

func (repo *GormRepository) createEmptyResultSet(payload interface{}) interface{} { // (pointerInterface interface{}, innerSlice interface{}, pointer reflect.Value) {

	payloadInfo := reflect.ValueOf(payload)
	payloadType := payloadInfo.Type()

	sliceType := reflect.SliceOf(payloadType)

	emptySlice := reflect.MakeSlice(sliceType, 0, 100000)

	pointerInstanceValueInfo := reflect.New(sliceType)

	pointerInstanceValueInfo.Elem().Set(emptySlice)

	pointerInstance := pointerInstanceValueInfo.Interface()

	return pointerInstance
}

func (repo *GormRepository) getIDValue(entity interface{}) uint64 {

	entityValue := reflect.ValueOf(entity)

	if entityValue.Kind() == reflect.Ptr {
		entityValue = entityValue.Elem()
	}

	return entityValue.FieldByName(ID_FIELD_NAME).Uint()
}
