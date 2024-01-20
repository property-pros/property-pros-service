package interop

// FinancialItem
type FinancialItem struct {
	Id           string `json:"id"`
	UserId       string `json:"userId"`
	ThirdPartyId string `json:"propertyId"`
	Token        string `json:"token"`
	Accounts     []FinancialAccount
	Transactions []FinancialTransaction
}

// Account represents the Acount message in the comment below
type FinancialAccount struct {
	AccountId    string            `json:"account_id"`
	Balances     *FinancialBalance `json:"balances"`
	Mask         string            `json:"mask"`
	Name         string            `json:"name"`
	OfficialName string            `json:"official_name"`
	Subtype      string            `json:"subtype"`
	Type         string            `json:"type"`
}

// Balance represents the Balance message in the comment below
type FinancialBalance struct {
	Available              float64 `json:"available"`
	Current                float64 `json:"current"`
	IsoCurrencyCode        string  `json:"iso_currency_code"`
	Limit                  float64 `json:"limit"`
	UnofficialCurrencyCode string  `json:"unofficial_currency_code"`
}

// Transaction represents the Transaction message in the comment below
type FinancialTransaction struct {
	AccountId              string       `json:"account_id"`
	Amount                 float64      `json:"amount"`
	IsoCurrencyCode        string       `json:"iso_currency_code"`
	UnofficialCurrencyCode string       `json:"unofficial_currency_code"`
	Category               []string     `json:"category"`
	CategoryId             string       `json:"category_id"`
	Date                   string       `json:"date"`
	Location               *FinancialLocation    `json:"location"`
	Name                   string       `json:"name"`
	PaymentMeta            *FinancialPaymentMeta `json:"payment_meta"`
	Pending                bool         `json:"pending"`
	PendingTransactionId   string       `json:"pending_transaction_id"`
	TransactionId          string       `json:"transaction_id"`
	TransactionType        string       `json:"transaction_type"`
}

// FinancialLocation represents the FinancialLocation message in the comment below
type FinancialLocation struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
}

// FinancialPaymentMeta represents the FinancialPaymentMeta message in the comment below
type FinancialPaymentMeta struct {
	ByOrderOf        string `json:"by_order_of"`
	Payee            string `json:"payee"`
	Payer            string `json:"payer"`
	PaymentMethod    string `json:"payment_method"`
	PaymentProcessor string `json:"payment_processor"`
	PpdId            string `json:"ppd_id"`
	Reason           string `json:"reason"`
	ReferenceNumber  string `json:"reference_number"`
}
