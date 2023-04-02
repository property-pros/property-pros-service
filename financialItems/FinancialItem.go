package data

import (
	"context"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

// FinancialItemModel
type FinancialItemModel struct {
	interop.FinancialItem
}

// GetPayload
func (m *FinancialItemModel) GetPayload() *interop.FinancialItem {
	return &m.FinancialItem
}

// SetPayload
func (m *FinancialItemModel) SetPayload(payload *interop.FinancialItem) {
	m.FinancialItem = *payload
}

// FinancialItemService
type FinancialItemService struct {
	FinancialItemGateway interfaces.IFinancialItemGateway
}

// SaveFinancialItem
func (s *FinancialItemService) SaveFinancialItem(ctx context.Context, payload *interop.FinancialItem) (*interop.FinancialItem, error) {
	
	model := &FinancialItemModel{}
	model.SetPayload(payload)

	err := s.FinancialItemGateway.SaveFinancialItem(ctx, model)

	if err != nil {
		return nil, err
	}

	return model.GetPayload(), nil
}

// GetFinancialItems
func (s *FinancialItemService) GetFinancialItems(ctx context.Context) ([]interfaces.IFinancialItemModel, error) {
	return s.FinancialItemGateway.GetFinancialItems(ctx)
}