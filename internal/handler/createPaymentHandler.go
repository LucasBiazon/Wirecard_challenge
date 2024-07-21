package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	dbPool *pgxpool.Pool
	logger *zap.Logger
)

func InitDB(pool *pgxpool.Pool, log *zap.Logger) {
	dbPool = pool
	logger = log
}

func CreatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	request := &PaymentRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Error("Invalid request payload", zap.Error(err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		logger.Error("Failed to start transaction", zap.Error(err))
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			logger.Error("Failed to rollback transaction", zap.Error(err))
		}
	}()

	var clientID string
	err = tx.QueryRow(ctx, `
			INSERT INTO client ()
			VALUES ()
			RETURNING id
		`).Scan(&clientID)
	if err != nil {
		logger.Error("Failed to insert client", zap.Error(err))
		http.Error(w, "Failed to insert client", http.StatusInternalServerError)
		return
	}

	var buyerID string
	err = tx.QueryRow(ctx, `
			INSERT INTO buyers (email, name, cpf)
			VALUES ($1, $2, $3)
			RETURNING id
		`, request.Buyer.Email, request.Buyer.Name, request.Buyer.CPF).Scan(&buyerID)
	if err != nil {
		logger.Error("Failed to insert buyer", zap.Error(err))
		http.Error(w, "Failed to insert buyer", http.StatusInternalServerError)
		return
	}

	var paymentID string
	if request.Payment.Type == "credit" {
		var creditID string
		err = tx.QueryRow(ctx, `
			INSERT INTO credit_card (number, exp_date, cvv, holder_name) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING id
		`, request.Payment.Card.Number, request.Payment.Card.ExpDate,
			request.Payment.Card.CVV, request.Payment.Card.HolderName).Scan(&creditID)
		if err != nil {
			logger.Error("Failed to insert credit card", zap.Error(err))
			http.Error(w, "Failed to insert credit card", http.StatusInternalServerError)
			return
		}

		paymentDate := time.Now()
		err = tx.QueryRow(ctx, `
			INSERT INTO payment (amount, type_payment, status_payment, payment_date, buyer_id, credit_card_id, client_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`, request.Payment.Amount, request.Payment.Type, "approved", paymentDate, buyerID, creditID, clientID).Scan(&paymentID)
		if err != nil {
			logger.Error("Failed to insert payment", zap.Error(err))
			http.Error(w, "Failed to insert payment", http.StatusInternalServerError)
			return
		}
	}

	if request.Payment.Type == "boleto" {
		err = tx.QueryRow(ctx, `
			INSERT INTO payment (amount, type_payment, status_payment,  payment_date, buyer_id, client_id)
			VALUES ($1, $2, $3, $4, $5, $6)`, request.Payment.Amount, request.Payment.Type, "pending", nil, buyerID, clientID).Scan(&paymentID)
		if err != nil {
			logger.Error("Failed to insert payment", zap.Error(err))
			http.Error(w, "Failed to insert payment", http.StatusInternalServerError)
			return
		}
	}
	var payment PaymentResponse
	err = tx.QueryRow(ctx, `
		SELECT id, amount, type_payment, status_payment, payment_date, buyer_id, client_id, credit_card_id
		FROM payment
		WHERE id = $1
	`, paymentID).Scan(&payment.ID, &payment.Amount, &payment.Type, &payment.Status, &payment.Date, &payment.BuyerID, &payment.ClientID, &payment.CreditCardID)
	if err != nil {
		logger.Error("Failed to query payment", zap.Error(err))
		http.Error(w, "Failed to query payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payment); err != nil {
		logger.Error("Failed to write response", zap.Error(err))
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
