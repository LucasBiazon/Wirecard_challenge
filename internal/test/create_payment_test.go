// create_payment_handler_test.go

package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lucasBiazon/wirecard/internal/handler"
	"go.uber.org/zap"
)

func TestCreatePaymentHandler(t *testing.T) {

	connString := "user=wirecard password=wirecard host=db port=5432 dbname=wirecard"

	// Criação do pool de conexões com o banco de dados
	dbPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Setup do logger
	logger, _ := zap.NewDevelopment()
	handler.InitDB(dbPool, logger)

	// Criação de uma solicitação de teste
	paymentRequest := handler.PaymentRequest{
		Client:  handler.Client{ID: "1"},
		Buyer:   handler.Buyer{ID: "1", Email: "buyer@example.com", Name: "Sample Buyer", CPF: "12345678909"},
		Payment: handler.Payment{Type: "credit", Amount: 100.0, Card: &handler.Card{Number: "4111111111111111", ExpDate: "12/23", CVV: "123", HolderName: "Sample Holder"}},
	}

	reqBody, err := json.Marshal(paymentRequest)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreatePaymentHandler(rr, req)

	// Verificar o código de status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verificar o corpo da resposta
	var response handler.PaymentResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Verificar a resposta esperada
	if response.Amount != 100.0 || response.Type != "credit" || response.Status != "approved" {
		t.Errorf("handler returned unexpected body: got %+v", response)
	}

	// Verificar se os IDs estão presentes e não são vazios
	if response.ID == "" || response.BuyerID == "" || response.ClientID == "" || response.CreditCardID == nil {
		t.Error("handler returned empty or invalid IDs")
	}
}
