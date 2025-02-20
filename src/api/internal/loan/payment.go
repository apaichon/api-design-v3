package loan

import "time"

type paymentService struct{}

func NewPaymentService() PaymentService {
    return &paymentService{}
}

func (s *paymentService) TransferFunds(fromAccount, toAccount string, amount float64) error {
    return nil
}

func (s *paymentService) ValidatePayment(paymentID string) error {
    return nil
}

func (s *paymentService) CalculateFine(dueDate time.Time, amount float64) float64 {
    return 0.0
} 