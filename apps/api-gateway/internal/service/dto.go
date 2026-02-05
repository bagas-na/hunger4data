package service

type ErrorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Auth
type RegisterRequest struct {
	Username string `json:"username" example:"bagas-na@github.com"`
	Password string `json:"password" example:"hunger4data"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	User    struct {
		Id       string `json:"id" example:"45dcf27b-29e4-46aa-a3aa-6e2b33ca86ae"`
		Username string `json:"username" example:"bagas-na@github.com"`
	} `json:"user"`
}

type LoginRequest struct {
	Username string `json:"username" example:"bagas-na@github.com"`
	Password string `json:"password" example:"hunger4data"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhdXRoLXNlcnZpY2UiLCJzdWIiOiI3M2E3ZjliMi04ODhlLTRmNzktYmU4YS0zNGU1NjExNmE4MmIiLCJhdWQiOlsiYm9va21zIl0sImV4cCI6MTc3MDI5OTM1MSwiaWF0IjoxNzcwMjg4NTUxfQ.kd1EJWYnAo9pw0R2_QH6Uxhjvmh10_AnPW3tbP55LC0"`
}

// Subscription
type CountryInfo struct {
	Id                        string  `json:"id,omitempty"`
	Name                      string  `json:"name,omitempty"`
	IpcPhase                  string  `json:"ipc_phase,omitempty" example:"3+"`
	PopulationInPhase         int64   `json:"population_in_phase,omitempty" example:"42000000"`
	PopulationFractionInPhase float64 `json:"population_fraction_in_phase,omitempty" example:"0.42"`
	LocationCode              string  `json:"location_code,omitempty" example:"COD"`
}

type GetCountriesResponse struct {
	Message string        `json:"message" example:"Success"`
	Data    []CountryInfo `json:"data"`
}

type SubcriptionInfo struct {
	Id          string `json:"id" example:"89aafe8b-e9aa-4daa-88f1-43e074ef0d6b"`
	UserId      string `json:"user_id" example:"73a7f9b2-888e-4f79-be8a-34e56116a82b"`
	CountryCode string `json:"country_code" example:"COF"`
}

type CreateSubcriptionRequest struct {
	CountryCode string `json:"country_code"`
}

type CreateSubcriptionResponse struct {
	Message string            `json:"message" example:"Subscription successful"`
	Data    []SubcriptionInfo `json:"data"`
}

type DeleteSubcriptionResponse struct {
	Message string `json:"message" example:"Delete successful"`
}

// Payments

type PaymentInfo struct {
	PaymentID       string `json:"paymentId"`
	CountryCode     string `json:"countryCode" example:"IDN"`
	TransactionType string `json:"transactionType" example:"donation"`
	Amount          int64  `json:"amount" example:"10000"`
	Currency        string `json:"currency" example:"USD"`
	Status          string `json:"status" example:"pending"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

type CreatePaymentRequest struct {
	CountryCode string `json:"country_code" example:"AFG"`
	Amount      int64  `json:"amount" example:"1000"`
	Currency    string `json:"currency" example:"usd"`
}

type CreatePaymentResponse struct {
	Message     string      `json:"message"`
	PaymentInfo PaymentInfo `json:"payment_info"`
	CheckoutUrl string      `json:"checkout_url"`
}

type GetPaymentCheckoutURLResponse struct {
	CheckoutUrl string `json:"checkout_url"`
}

type ListPaymentsResponse []PaymentInfo
