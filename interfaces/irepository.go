package interfaces

type IRepository[T any] interface {
	Get(string) (*T, error)
	Save(payload *T) (*T, error)
	FindOne(*T) (*T, error)
	Query(*T) []*T
	Delete(*T) (*T, error)
}
