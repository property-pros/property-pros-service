package interfaces

import (
	"context"
)

type IBaseModel[T any] interface {
	GetPayload() *T
	GetContext() context.Context
	SetPayload(*T)
	SetContext(context.Context)
	GetId() string
}
