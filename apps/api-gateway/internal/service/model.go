package service

type PaymentInfo struct {
	PaymentID       string `json:"payment_id"`
	CountryCode     string `json:"country_code"`
	TransactionType string `json:"transaction_type"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type CreatePaymentResponse struct {
	Message     string      `json:"message"`
	PaymentInfo PaymentInfo `json:"payment_info"`
	CheckoutUrl string      `json:"checkout_url"`
}

type ListPaymentsResponse []PaymentInfo
