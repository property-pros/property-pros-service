package bootstrap

import (
	ctx "context"

	"github.com/vireocloud/property-pros-service/interfaces"
)

func NewModel[T, M interfaces.IBaseModel[T]](context ctx.Context, payload *T, initializer func() (M, error)) (M, error) {
	model, err := initializer()

	if err != nil {
		return model, err
	}

	model.SetContext(context)
	model.SetPayload(payload)

	return model, nil
}
