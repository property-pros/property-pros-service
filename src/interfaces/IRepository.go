package interfaces

type IRepository interface {
	Get(uint)
	Save(payload interface{}, query interface{}) (interface{}, error)
	FindOne(interface{}) (interface{}, error)
	Query(interface{}) []interface{}
	Delete(interface{}) (interface{}, error)
	Initialize() error
	GetQueryAndPayload(payload interface{}) (newPayload interface{}, query interface{})
}
