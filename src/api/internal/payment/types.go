package payment

type Payment struct {
	PaymentID     string  `json:"payment_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	PaymentDate   string  `json:"payment_date"`
	PayTo         string  `json:"pay_to"`
	Note          string  `json:"note"`
	Status        string  `json:"status"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type CreditCard struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	CVV        string `json:"cvv"`
}

type CreditCardPayment struct {
	PaymentID    string     `json:"payment_id"`
	CreditCardInfo CreditCard `json:"credit_card_info"`
	Amount         float64    `json:"amount"`
	PayTo          string     `json:"pay_to"`
	Note           string     `json:"note"`
	Status         string     `json:"status"`
	CreatedAt      string     `json:"created_at"`
	UpdatedAt      string     `json:"updated_at"`
}
