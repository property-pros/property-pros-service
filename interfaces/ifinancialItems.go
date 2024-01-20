package interfaces

import (
	"context"

	"github.com/vireocloud/property-pros-service/interop"
)

// IFinancialItemRepository
type IFinancialItemGateway interface {
	SaveFinancialItem(ctx context.Context, model IFinancialItemModel) error
	GetFinancialItems(ctx context.Context) ([]IFinancialItemModel, error)
}

// IFinancialItemModel
type IFinancialItemModel interface {
	GetPayload() *interop.FinancialItem
	SetPayload(payload *interop.FinancialItem)
}

// IFinancialService
type IFinancialService interface {
	SaveFinancialItem(ctx context.Context, payload *interop.FinancialItem) (*interop.FinancialItem, error)
	GetFinancialItems(ctx context.Context) ([]IFinancialItemModel, error)
}