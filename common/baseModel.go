package common

import (
	"context"
)

type BaseModel[T any] struct {
	Payload *T
	Context context.Context
}

func (model *BaseModel[T]) SetPayload(payload *T) {
	model.Payload = payload
}

func (model *BaseModel[T]) SetContext(ctx context.Context) {
	model.Context = ctx
}

func (model *BaseModel[T]) GetPayload() *T {
	return model.Payload
}

func (model *BaseModel[T]) GetContext() context.Context {
	return model.Context
}

// func (model *BaseModel[T]) GetId() string {
// 	return PT(model.Payload).GetId()
// }
