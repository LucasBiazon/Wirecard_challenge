package handler

type Card struct {
	Number     string `json:"number"`
	ExpDate    string `json:"exp_date"`
	CVV        string `json:"cvv"`
	HolderName string `json:"holder_name"`
}

type Payment struct {
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
	Card   *Card   `json:"card"`
}

type Buyer struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	CPF   string `json:"cpf"`
}

type Client struct {
	ID string `json:"id"`
}

type PaymentRequest struct {
	Client  Client  `json:"client"`
	Buyer   Buyer   `json:"buyer"`
	Payment Payment `json:"payment"`
}
