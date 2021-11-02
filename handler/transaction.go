package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/handler/middleware"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/usecase"
)

type TransactionHandler struct {
	usecase.TransactionUseCase
}

func NewTransactionHandler(u usecase.TransactionUseCase) TransactionHandler {
	return TransactionHandler{u}
}

func (s *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	var body entity.TransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}
	
	body.User_Id = claims.UserID
	if valid, msg := helper.Validate(body); !valid {
		badRequest(w, msg)
	}
	if err := s.TransactionUseCase.MakeTransaction(ctx, body); err != nil {
		log.Println(err)
		internalServerError(w)
		return
	}
	resp, _ := json.Marshal(commonResponse{
		Message: "resources has successfully created",
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
	}

	transactions, err := s.TransactionUseCase.GetUserTransactions(ctx, claims.UserID)

	if err != nil {
		if err.Error() == "object storage service unavailable" {
			serviceUnavailable(w, "error: "+err.Error())
			return
		}
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

func (s *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload *entity.Transaction `json:"payload"`
	}

	ctx := r.Context()
	transactionID, _ := strconv.Atoi(chi.URLParam(r, "transactionID"))

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
	transactionID, _ := strconv.Atoi(chi.URLParam(r, "transactionID"))

	var data map[string]interface{}

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
