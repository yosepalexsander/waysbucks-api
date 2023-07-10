package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/middleware"
	"github.com/yosepalexsander/waysbucks-api/usecase"
)

type AddressHandler struct {
	AddressUseCase usecase.AddressUseCase
}

func (s *AddressHandler) FindUserAddresses(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload []entity.Address `json:"payload"`
	}

	ctx := r.Context()
	claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	address, err := s.AddressUseCase.FindUserAddresses(ctx, claims.UserID)
	if err != nil {
		internalServerError(w)
		return
	}

	responseStruct := response{
		commonResponse: commonResponse{
			Message: "resource has successfully get",
		},
		Payload: address,
	}
	resp, _ := json.Marshal(responseStruct)
	responseOK(w, resp)
}

func (s *AddressHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload *entity.Address `json:"payload"`
	}

	addressID := chi.URLParam(r, "addressID")
	address, err := s.AddressUseCase.GetAddress(r.Context(), addressID)

	if err != nil {
		internalServerError(w)
		return
	}

	responseStruct := response{
		commonResponse: commonResponse{
			Message: "resource has successfully get",
		},
		Payload: address,
	}

	resp, _ := json.Marshal(responseStruct)
	responseOK(w, resp)
}

func (s *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var body entity.AddressRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}

	if isValid, msg := helper.Validate(body); !isValid {
		badRequest(w, msg)
		return
	}

	if err := s.AddressUseCase.CreateNewAddress(ctx, claims.UserID, body); err != nil {
		internalServerError(w)
		return
	}

	resp, _ := json.Marshal(commonResponse{
		Message: "resource successfully created",
	})
	responseOK(w, resp)
}

func (s *AddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	addressID := chi.URLParam(r, "addressID")

	address, err := s.AddressUseCase.GetAddress(ctx, addressID)

	if err != nil {
		notFound(w)
		return
	}

	if claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims); ok {
		if claims.UserID != address.UserId {
			forbidden(w)
			return
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}

	if err := s.AddressUseCase.UpdateAddress(ctx, addressID, body); err != nil {
		internalServerError(w)
		return
	}

	resp, _ := json.Marshal(commonResponse{
		Message: "resource successfully updated",
	})
	responseOK(w, resp)
}

func (s *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	addressID := chi.URLParam(r, "addressID")

	if claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims); ok {
		if err := s.AddressUseCase.DeleteAddress(ctx, addressID, claims.UserID); err != nil {
			if err.Error() == "no rows affected" {
				notFound(w)
				return
			}
			internalServerError(w)
			return
		}

		resBody, _ := json.Marshal(commonResponse{
			Message: "resource successfully deleted",
		})
		responseOK(w, resBody)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
}
