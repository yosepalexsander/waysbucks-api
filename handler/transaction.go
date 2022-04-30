package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/middleware"
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
	"github.com/yosepalexsander/waysbucks-api/usecase"
)

type TransactionHandler struct {
	usecase.TransactionUseCase
}

func NewTransactionHandler(u usecase.TransactionUseCase) TransactionHandler {
	return TransactionHandler{u}
}

func (s *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	type ResponsePayload struct {
		Token       string `json:"token"`
		RedirectURL string `json:"redirect_url"`
	}
	type response struct {
		commonResponse
		Payload ResponsePayload `json:"payload"`
	}

	ctx := r.Context()
	claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	body := entity.TransactionRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}

	body.UserId = claims.UserID
	if valid, msg := helper.Validate(body); !valid {
		badRequest(w, msg)
	}

	createdTransaction, err := s.TransactionUseCase.MakeTransaction(ctx, body)
	if err != nil {
		internalServerError(w)
		return
	}

	snapRes := thirdparty.CreateTransaction(createdTransaction)

	resp, _ := json.Marshal(response{
		commonResponse: commonResponse{
			Message: "resources has successfully created",
		},
		Payload: ResponsePayload{
			Token:       snapRes.Token,
			RedirectURL: snapRes.RedirectURL,
		},
	})

	responseOK(w, resp)
}

func (s *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload []entity.Transaction `json:"payload"`
	}

	ctx := r.Context()
	_, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	transactions, err := s.TransactionUseCase.GetTransactions(ctx)
	if err != nil {
		internalServerError(w)
		return
	}

	resp, _ := json.Marshal(response{
		commonResponse: commonResponse{
			Message: "resources has successfully get",
		},
		Payload: transactions,
	})

	responseOK(w, resp)
}

func (s *TransactionHandler) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload []entity.Transaction `json:"payload"`
	}

	ctx := r.Context()
	claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	transactions, err := s.TransactionUseCase.GetUserTransactions(ctx, claims.UserID)
	if err != nil {
		switch err {
		case thirdparty.ErrServiceUnavailable:
			serviceUnavailable(w, "error: cloudinary service unavailable")
		case sql.ErrNoRows:
			notFound(w)
		default:
			internalServerError(w)
		}
		return
	}

	resp, _ := json.Marshal(response{
		commonResponse: commonResponse{
			Message: "resources has successfully get",
		},
		Payload: transactions,
	})

	responseOK(w, resp)
}

func (s *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload *entity.Transaction `json:"payload"`
	}

	ctx := r.Context()
	transactionID := chi.URLParam(r, "transactionID")

	transaction, err := s.TransactionUseCase.GetDetailTransaction(ctx, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			notFound(w)
			return
		}
		internalServerError(w)
		return
	}

	resp, _ := json.Marshal(response{
		commonResponse: commonResponse{
			Message: "resource has successfully get",
		},
		Payload: transaction,
	})

	responseOK(w, resp)
}

func (s *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := chi.URLParam(r, "transactionID")

	data := make(map[string]interface{})

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		badRequest(w, "invalid request")
		return
	}

	if err := s.TransactionUseCase.UpdateTransaction(r.Context(), transactionID, data); err != nil {
		internalServerError(w)
		return
	}

	resp, _ := json.Marshal(commonResponse{
		Message: "resource has successfully updated",
	})

	responseOK(w, resp)
}

// Catch notification from midtrans request POST after
func (s *TransactionHandler) PaymentNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	transaction, err := thirdparty.ParseTransactionResponse(r.Body)
	if err != nil {
		internalServerError(w)
		return
	}

	updateStatus := func(status string) {
		data := make(map[string]interface{})
		data["status"] = status

		if err := s.TransactionUseCase.UpdateTransaction(ctx, transaction.TransactionID, data); err != nil {
			internalServerError(w)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

	if transaction.TransactionStatus == "capture" {
		if transaction.FraudStatus == "challange" {
			updateStatus("pending")
		} else if transaction.FraudStatus == "accept" {
			updateStatus("success")
		}
	} else if transaction.TransactionStatus == "settlement" {
		updateStatus("success")
	} else if transaction.TransactionStatus == "cancel" || transaction.TransactionStatus == "deny" || transaction.TransactionStatus == "expire" {
		updateStatus("failure")
	} else if transaction.TransactionStatus == "pending" {
		updateStatus("pending")
	}
}
