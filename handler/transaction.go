package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yosepalexsander/waysbucks-api/usecase"
)


type TransactionHandler struct {
	usecase.TransactionUseCase
}

func NewTransactionHandler(u usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{u}
}

func (s *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request)  {
	type request struct {
		ProductId int `json:"product_id"`
		Qty int `json:"qty"`
		Price int `json:"price"`
		ToppingIds []int `json:"topping_id"`
	}

	var body []request

	if err:= json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}
}