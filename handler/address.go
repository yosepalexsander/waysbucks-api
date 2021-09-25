package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/handler/middleware"
	"github.com/yosepalexsander/waysbucks-api/persistance"
)
type AddressServer struct {
	Repo persistance.AddressRepository
}

func (s *AddressServer) GetAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(middleware.TokenCtxKey).(*middleware.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	
	address, err := s.Repo.FindUserAddress(ctx, claims.UserID);
	if  err != nil {
		internalServerError(w)
		return
	}

	responseStruct := struct{
		CommonResponse
		Payload *[]entity.UserAddress `json:"payload"`
	} {
		CommonResponse: CommonResponse{
			Message: "resource has successfully get",
		},
		Payload: address,
	}
	resp, _ := json.Marshal(responseStruct)
	responseOK(w, resp)
}
func (s *AddressServer) CreateAddress(w http.ResponseWriter, r *http.Request)  {
	ctx := r.Context()
	claims, ok := ctx.Value(middleware.TokenCtxKey).(*middleware.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var body entity.UserAddress
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
    return
	}
	body.UserId = claims.UserID

	if err := s.Repo.SaveAddress(ctx, body); err != nil {
		internalServerError(w)
		return
	}

	resp, _ := json.Marshal(CommonResponse{
		Message: "resource successfully created",
	})
	responseOK(w, resp)
}

func (s *AddressServer) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	addressID, _:= strconv.Atoi(chi.URLParam(r, "addressID"))

	address, err := s.Repo.FindAddress(ctx, addressID)

	if err != nil {
		notFound(w)
		return
	}
	
	if claims, ok := ctx.Value(middleware.TokenCtxKey).(*middleware.MyClaims); ok {
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

	if err := s.Repo.UpdateAddress(ctx, addressID, body); err != nil {
		internalServerError(w)
		return
	}

	resp, _ := json.Marshal(CommonResponse{
		Message: "resource successfully updated",
	})
	responseOK(w, resp)
}

func (s *AddressServer) DeleteAddress(w http.ResponseWriter, r *http.Request)  {
	ctx := r.Context()
	addressID, _ := strconv.Atoi(chi.URLParam(r, "addressID"))
	
	if claims, ok := ctx.Value(middleware.TokenCtxKey).(*middleware.MyClaims); ok {
		if err := s.Repo.DeleteAddress(ctx, addressID, claims.UserID); err != nil {
			if err.Error() == "no rows affected" {
				notFound(w)
				return
			}
			internalServerError(w)
			return
		}
		
		resBody, _ := json.Marshal(CommonResponse{
			Message:  "resource successfully deleted",
		})
		responseOK(w, resBody)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
}