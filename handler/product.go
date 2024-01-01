package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
	"github.com/yosepalexsander/waysbucks-api/usecase"
)

type ProductHandler struct {
	ProductUseCase usecase.ProductUseCase
}

func NewProductHandler(u usecase.ProductUseCase) ProductHandler {
	return ProductHandler{u}
}

func (s *ProductHandler) FindProducts(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload []entity.Product `json:"payload"`
	}

	queries := r.URL.Query()
	products, err := s.ProductUseCase.FindProducts(r.Context(), queries)

	if err != nil {
		internalServerError(w)
		return
	}

	responseStruct := response{
		commonResponse: commonResponse{
			Message: "resource has successfully get",
		},
		Payload: products,
	}

	resp, _ := json.Marshal(responseStruct)
	responseOK(w, resp)
}

func (s *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload *entity.Product `json:"payload"`
	}

	ctx := r.Context()
	productID, _ := strconv.Atoi(chi.URLParam(r, "productID"))

	product, err := s.ProductUseCase.GetProduct(ctx, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			notFound(w)
			return
		}
		internalServerError(w)
		return
	}

	responseStruct := response{
		commonResponse: commonResponse{
			Message: "resource has successfully get",
		},
		Payload: product,
	}

	resp, _ := json.Marshal(responseStruct)

	responseOK(w, resp)
}

func (s *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := entity.ProductRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}

	if valid, msg := helper.Validate(body); !valid {
		badRequest(w, msg)
		return
	}

	if err := s.ProductUseCase.CreateProduct(ctx, body); err != nil {
		internalServerError(w)
		return
	}

	resp, _ := json.Marshal(commonResponse{
		Message: "resource has successfully created",
	})

	responseOK(w, resp)
}

func (s *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID, _ := strconv.Atoi(chi.URLParam(r, "productID"))
	body := make(map[string]interface{})

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request body")
		return
	}

	if err := s.ProductUseCase.UpdateProduct(ctx, productID, body); err != nil {
		internalServerError(w)
		return
	}

	resBody, _ := json.Marshal(commonResponse{
		Message: "resource has successfully updated",
	})

	responseOK(w, resBody)
}

func (s *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID, _ := strconv.Atoi(chi.URLParam(r, "productID"))

	if err := s.ProductUseCase.DeleteProduct(ctx, productID); err != nil {
		internalServerError(w)
		return
	}

	resBody, _ := json.Marshal(commonResponse{
		Message: "resource has successfully deleted",
	})

	responseOK(w, resBody)
}

func (s *ProductHandler) FindToppings(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload []entity.ProductTopping `json:"payload"`
	}

	toppings, err := s.ProductUseCase.FindToppings(r.Context())
	if err != nil {
		switch err {
		case thirdparty.ErrServiceUnavailable:
			serviceUnavailable(w, "error: cloudinary service unavailable")
		default:
			internalServerError(w)
		}

		return
	}

	responseStruct := response{
		commonResponse: commonResponse{
			Message: "resource has successfully get",
		},
		Payload: toppings,
	}
	resBody, _ := json.Marshal(responseStruct)

	responseOK(w, resBody)
}

func (s *ProductHandler) CreateTopping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body := entity.ProductToppingRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request body")
		return
	}

	if valid, msg := helper.Validate(body); !valid {
		badRequest(w, msg)
		return
	}

	if err := s.ProductUseCase.CreateTopping(ctx, body); err != nil {
		internalServerError(w)
		return
	}

	resBody, _ := json.Marshal(commonResponse{
		Message: "resource has successfully created",
	})

	responseOK(w, resBody)
}

func (s *ProductHandler) UpdateTopping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	toppingID, _ := strconv.Atoi(chi.URLParam(r, "toppingID"))
	body := make(map[string]interface{})

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request body")
		return
	}

	if err := s.ProductUseCase.UpdateTopping(ctx, toppingID, body); err != nil {
		internalServerError(w)
		return
	}

	resBody, _ := json.Marshal(commonResponse{
		Message: "resource has successfully updated",
	})

	responseOK(w, resBody)
}

func (s *ProductHandler) DeleteTopping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	toppingID, _ := strconv.Atoi(chi.URLParam(r, "toppingID"))

	if err := s.ProductUseCase.DeleteTopping(ctx, toppingID); err != nil {
		internalServerError(w)
		return
	}

	resBody, _ := json.Marshal(commonResponse{
		Message: "resource has successfully deleted",
	})

	responseOK(w, resBody)
}
