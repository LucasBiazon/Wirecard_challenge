package handler

import "time"

type PaymentResponse struct {
	ID           string    `json:"id"`
	Amount       float64   `json:"amount"`
	Type         string    `json:"type_payment"`
	Status       string    `json:"status_payment"`
	Date         time.Time `json:"payment_date"`
	BuyerID      string    `json:"buyer_id"`
	ClientID     string    `json:"client_id"`
	CreditCardID *string   `json:"credit_card_id,omitempty"` // Pode ser nulo para pagamentos com boleto
}
