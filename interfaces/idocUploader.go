package interfaces

import "context"

type IDocUploader interface {
	PutObject(ctx context.Context, content []byte) (string, error)
}
