package interfaces

type IEntityGateway interface {
	SaveEntities(entities []interface{}) []interface{}
	SaveEntity(entity interface{}) interface{}
	DeleteEntity(entity interface{}) interface{}
	SaveEntityWithMatch(entity interface{}, idQuery interface{}) interface{}
	ListEntities(query interface{}) []interface{}
}
