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
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
	"github.com/yosepalexsander/waysbucks-api/usecase"
)

type CartHandler struct {
	CartUseCase usecase.CartUseCase
}

func (s *CartHandler) GetUserCarts(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload []entity.Cart `json:"payload"`
	}
	ctx := r.Context()

	if claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims); ok {
		carts, err := s.CartUseCase.GetUserCarts(ctx, claims.UserID)
		if err != nil {
			if err == thirdparty.ErrServiceUnavailable {
				serviceUnavailable(w, "error: cloudinary service unavailable")
				return
			}
			if err == sql.ErrNoRows {
				notFound(w)
				return
			}
			internalServerError(w)
			return
		}

		resp, _ := json.Marshal(response{
			commonResponse: commonResponse{
				Message: "resources successfully get",
			},
			Payload: carts,
		})

		responseOK(w, resp)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
}

func (s *CartHandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var body entity.Cart
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}

	body.User_Id = claims.UserID

	if valid, msg := helper.Validate(body); !valid {
		badRequest(w, msg)
		return
	}

	err := s.CartUseCase.SaveToCart(ctx, body)
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
	cartId, _ := strconv.Atoi(chi.URLParam(r, "cartID"))
	claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}

	if err := s.CartUseCase.UpdateCart(ctx, cartId, claims.UserID, body); err != nil {
		internalServerError(w)
		log.Println(err)
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
