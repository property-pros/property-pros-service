package interfaces

type IFactory[T any] interface {
	New() *T
}
