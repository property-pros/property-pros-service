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

func (c *FinanceController) GetFinancialItems(ctx context.Context, req *interop.GetFinancialItemsRequest) (response *interop.GetFinancialItemsResponse, errResult error) {

	response = &interop.GetFinancialItemsResponse{}

	query := &interop.FinancialItem{
		UserId: req.GetUserId(),
	}

	// get financial item from repo
	financialItems := c.financialRepo.Query(query)

	// map FinancialData to response
	response.FinancialData = []*interop.FinancialData{}

	for _, financialItem := range financialItems {

		// map FinancialItem to FinancialData
		financialData := &interop.FinancialData{
			ItemId:       financialItem.ThirdPartyId,
			AccessToken:  financialItem.Token,
			Accounts:     []*interop.Account{},
			Transactions: []*interop.Transaction{},
		}

		// map FinancialAccount to FinancialData
		for _, account := range financialItem.Accounts {

			// map FinancialBalance to FinancialAccount
			balances := &interop.Balance{
				Available:              account.Balances.Available,
				Current:                account.Balances.Current,
				Limit:                  account.Balances.Limit,
				UnofficialCurrencyCode: account.Balances.UnofficialCurrencyCode,
				IsoCurrencyCode:        account.Balances.IsoCurrencyCode,
			}

			financialData.Accounts = append(financialData.Accounts, &interop.Account{
				AccountId:    account.AccountId,
				Type:         account.Type,
				Name:         account.Name,
				Subtype:      account.Subtype,
				Mask:         account.Mask,
				OfficialName: account.OfficialName,
				Balances:     balances,
			})
		}

		// map FinancialTransaction to FinancialData
		for _, transaction := range financialItem.Transactions {

			// map Location to FinancialTransaction
			location := &interop.Location{
				Address: transaction.Location.Address,
				City:    transaction.Location.City,
				Region: transaction.Location.State,
				PostalCode: transaction.Location.Zip,
			}

			// map PaymentMeta to FinancialTransaction
			paymentMeta := &interop.PaymentMeta{
				Payee:            transaction.PaymentMeta.Payee,
				Payer:            transaction.PaymentMeta.Payer,
				PaymentMethod:    transaction.PaymentMeta.PaymentMethod,
				PaymentProcessor: transaction.PaymentMeta.PaymentProcessor,
				PpdId:            transaction.PaymentMeta.PpdId,
				ReferenceNumber:  transaction.PaymentMeta.ReferenceNumber,
			}

			// map FinancialTransaction to FinancialData
			transactionPayload := &interop.Transaction{
				AccountId:    transaction.AccountId,
				TransactionId: transaction.TransactionId,
				Amount:       transaction.Amount,
				Category:     transaction.Category,
				CategoryId:   transaction.CategoryId,
				Date:         transaction.Date,
				Name:         transaction.Name,
				Pending:      transaction.Pending,
				PendingTransactionId: transaction.PendingTransactionId,
				TransactionType:      transaction.TransactionType,
				Location:             location,
				PaymentMeta:          paymentMeta,				
			}

			financialData.Transactions = append(financialData.Transactions, transactionPayload)
		}

		response.FinancialData = append(response.FinancialData, financialData)
	}

	return response, nil
}

func (c *FinanceController) SaveFinancialItem(ctx context.Context, req *interop.SaveFinancialItemRequest) (response *interop.SaveFinancialItemResponse, errResult error) {

	response = &interop.SaveFinancialItemResponse{}

	// map request to FinancialItem
	payload := &interop.FinancialItem{
		UserId:       req.GetPayload().GetUserId(),
		ThirdPartyId: req.GetPayload().GetItemId(),
		Token:        req.GetPayload().GetAccessToken(),
	}

	// loop through req.Accounts and map to payload.Accounts
	for _, account := range req.GetPayload().GetAccounts() {

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
	for _, transaction := range req.GetPayload().GetTransactions() {

		// map request location to interop.FinancialLocation
		location := &interop.FinancialLocation{
			Address: transaction.GetLocation().GetAddress(),
			City:    transaction.GetLocation().GetCity(),
			State:   transaction.GetLocation().GetRegion(),
			Zip:     transaction.GetLocation().GetPostalCode(),
		}

		paymentMeta := &interop.FinancialPaymentMeta{
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
