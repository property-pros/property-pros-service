package controllers

import (
	"context"

	"github.com/vireocloud/property-pros-service/interfaces"
	"github.com/vireocloud/property-pros-service/interop"
)

// FinanceController
type FinanceController struct {
	interop.UnimplementedFinanceServiceServer

	financialRepo interfaces.IRepository[interop.FinancialItem]
}

// GetFinancialItem

func (c *FinanceController) GetFinancialItem(ctx context.Context, req *interop.GetFinancialItemRequest) (response *interop.GetFinancialItemResponse, errResult error) {

	response = &interop.GetFinancialItemResponse{}

	// get financial item from repo
	financialItem, err := c.financialRepo.Get(ctx, req.GetUserId())

	if err != nil {
		return nil, err
	}

	// map FinancialData to response
	response.FinancialData = &interop.FinancialData{
		AccessToken:  financialItem.Token,
		ItemId:       financialItem.ThirdPartyId,
		UserId:       financialItem.UserId,
		Accounts:     financialItem.Accounts,
		Transactions: financialItem.Transactions,
	}

	return response, nil
}

func (c *FinanceController) SaveFinancialItem(ctx context.Context, req *interop.SaveFinancialItemRequest) (response *interop.SaveFinancialItemResponse, errResult error) {

	response = &interop.SaveFinancialItemResponse{}

	// map request to FinancialItem
	payload := &interop.FinancialItem{
		UserId:       req.GetUserId(),
		ThirdPartyId: req.GetItemId(),
		Token:        req.GetAccessToken(),
	}

	// loop through req.Accounts and map to payload.Accounts
	for _, account := range req.GetAccounts() {

		// map request balances to interop.FinancialBalance
		balances := &interop.FinancialBalance{
			Available:              account.GetBalances().GetAvailable(),
			Current:                account.GetBalances().GetCurrent(),
			Limit:                  account.GetBalances().GetLimit(),
			UnofficialCurrencyCode: account.GetBalances().GetUnofficialCurrencyCode(),
			IsoCurrencyCode:        account.GetBalances().GetIsoCurrencyCode(),
		}

		payload.Accounts = append(payload.Accounts, interop.FinancialAccount{
			AccountId:    account.GetAccountId(),
			Type:         account.GetType(),
			Name:         account.GetName(),
			Subtype:      account.GetSubtype(),
			Mask:         account.GetMask(),
			OfficialName: account.GetOfficialName(),
			Balances:     balances,
		})
	}

	// loop through req.Transactions and map to payload.Transactions
	for _, transaction := range req.GetTransactions() {

		// map request location to interop.FinancialLocation
		location := &interop.Location{
			Address: transaction.GetLocation().GetAddress(),
			City:    transaction.GetLocation().GetCity(),
			State:   transaction.GetLocation().GetRegion(),
			Zip:     transaction.GetLocation().GetPostalCode(),
		}

		paymentMeta := &interop.PaymentMeta{
			ByOrderOf: transaction.GetPaymentMeta().GetByOrderOf(),
		}

		// map request transaction to interop.FinancialTransaction
		payload.Transactions = append(payload.Transactions, interop.FinancialTransaction{
			AccountId:              transaction.GetAccountId(),
			Amount:                 transaction.GetAmount(),
			IsoCurrencyCode:        transaction.GetIsoCurrencyCode(),
			UnofficialCurrencyCode: transaction.GetUnofficialCurrencyCode(),
			Category:               transaction.GetCategory(),
			CategoryId:             transaction.GetCategoryId(),
			Date:                   transaction.GetDate(),
			Location:               location,
			Name:                   transaction.GetName(),
			PaymentMeta:            paymentMeta,
			Pending:                transaction.GetPending(),
			PendingTransactionId:   transaction.GetPendingTransactionId(),
			TransactionId:          transaction.GetTransactionId(),
			TransactionType:        transaction.GetTransactionType(),
		})
	}

	_, err := c.financialRepo.Save(payload)

	if err != nil {
		return response, err
	}

	return response, nil
}
