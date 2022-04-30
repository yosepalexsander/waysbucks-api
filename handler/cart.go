package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/middleware"
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
	"github.com/yosepalexsander/waysbucks-api/usecase"
)

type CartHandler struct {
	CartUseCase usecase.CartUseCase
}

func NewCartHandler(u usecase.CartUseCase) CartHandler {
	return CartHandler{u}
}

func (s *CartHandler) GetCarts(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload []entity.Cart `json:"payload"`
	}

	ctx := r.Context()
	claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	carts, err := s.CartUseCase.GetCarts(ctx, claims.UserID)
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
			Message: "resources successfully get",
		},
		Payload: carts,
	})

	responseOK(w, resp)
}

func (s *CartHandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var body entity.CartRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}

	if valid, msg := helper.Validate(body); !valid {
		badRequest(w, msg)
		return
	}

	err := s.CartUseCase.SaveCart(ctx, body, claims.UserID)
	if err != nil {
		internalServerError(w)
		return
	}

	resBody, _ := json.Marshal(commonResponse{
		Message: "resource has successfully created",
	})

	responseOK(w, resBody)
}

func (s *CartHandler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cartID, _ := strconv.Atoi(chi.URLParam(r, "cartID"))
	claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	body := make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}

	if err := s.CartUseCase.UpdateCart(ctx, cartID, claims.UserID, body); err != nil {
		internalServerError(w)
		return
	}

	resBody, _ := json.Marshal(commonResponse{
		Message: "resource has successfully updated",
	})

	responseOK(w, resBody)
}

func (s CartHandler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cartID, _ := strconv.Atoi(chi.URLParam(r, "cartID"))
	claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := s.CartUseCase.DeleteCart(ctx, cartID, claims.UserID); err != nil {
		internalServerError(w)
		return
	}

	resBody, _ := json.Marshal(commonResponse{
		Message: "resource has successfully deleted",
	})

	responseOK(w, resBody)
}
