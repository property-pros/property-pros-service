package interfaces

type IRepository[T any] interface {
	Get(uint)
	Save(payload *T) (*T, error)
	FindOne(*T) (*T, error)
	Query(*T) []*T
	Delete(*T) (*T, error)
}
