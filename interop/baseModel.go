package interop

import "context"

type BaseModel[T any] struct {
	Payload *T
	Context context.Context
}
