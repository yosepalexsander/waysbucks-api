package handler

import (
	"database/sql"
	"encoding/json"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/middleware"
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
	"github.com/yosepalexsander/waysbucks-api/usecase"
)

type ProductHandler struct {
	ProductUseCase usecase.ProductUseCase
}

func NewProductHandler(u usecase.ProductUseCase) ProductHandler {
	return ProductHandler{u}
}

func (s *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload []entity.Product `json:"payload"`
	}

	queries := r.URL.Query()
	products, err := s.ProductUseCase.GetProducts(r.Context(), queries)

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

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		badRequest(w, "maximum upload size is 5 MB")
		return
	}

	file, header, fileErr := r.FormFile("image")
	if fileErr != nil {
		badRequest(w, fileErr.Error())
		return
	}
	defer file.Close()

	if err := helper.ValidateImageFile(header.Header.Get("Content-Type")); err != nil {
		badRequest(w, "upload only for image")
		return
	}

	body := entity.ProductRequest{}

	if err := schema.NewDecoder().Decode(&body, r.MultipartForm.Value); err != nil {
		badRequest(w, err.Error())
		return
	}

	filename, err := thirdparty.UploadFile(ctx, file, header.Filename)
	if err != nil {
		internalServerError(w)
		return
	}

	body.Image = filename
	if err := s.ProductUseCase.CreateProduct(ctx, body); err != nil {
		thirdparty.RemoveFile(ctx, filename)
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
	mediatype, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	body := make(map[string]interface{})

	if mediatype == "multipart/form-data" {
		if err := r.ParseMultipartForm(5 << 20); err != nil {
			badRequest(w, "maximum upload size is 5 MB")
			return
		}

		for k, v := range r.MultipartForm.Value {
			body[k] = v[0]
		}

		file, header, fileErr := r.FormFile("image")
		if fileErr != nil {
			badRequest(w, fileErr.Error())
			return
		}
		defer file.Close()

		product, err := s.ProductUseCase.GetProduct(ctx, productID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				notFound(w)
			default:
				internalServerError(w)
			}

			return
		}

		filename, err := thirdparty.UpdateImage(file, product.Image, header.Filename)
		if err != nil {
			internalServerError(w)
			return
		}

		body["image"] = filename

		if err := s.ProductUseCase.UpdateProduct(ctx, productID, body); err != nil {
			internalServerError(w)
			return
		}

		resBody, _ := json.Marshal(commonResponse{
			Message: "resource has successfully updated",
		})
		responseOK(w, resBody)
		return
	}

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
	_, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	product, err := s.ProductUseCase.GetProduct(ctx, productID)
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

	if err := s.ProductUseCase.DeleteProduct(ctx, productID); err != nil {
		internalServerError(w)
		return
	}

	if err := thirdparty.RemoveFile(ctx, product.Image); err != nil {
		internalServerError(w)
		return
	}

	resBody, _ := json.Marshal(commonResponse{
		Message: "resource has successfully deleted",
	})

	responseOK(w, resBody)
}

func (s *ProductHandler) GetToppings(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload []entity.ProductTopping `json:"payload"`
	}

	toppings, err := s.ProductUseCase.GetToppings(r.Context())
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

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		badRequest(w, "maximum upload size is 5 MB")
		return
	}

	file, header, fileErr := r.FormFile("image")
	if fileErr != nil {
		badRequest(w, fileErr.Error())
		return
	}

	defer file.Close()

	if err := helper.ValidateImageFile(header.Header.Get("Content-Type")); err != nil {
		badRequest(w, "upload only for image")
	}

	filename := strings.Split(header.Filename, ".")[0] + "-" + helper.RandString(15)
	body := entity.ProductToppingRequest{}

	if err := schema.NewDecoder().Decode(&body, r.MultipartForm.Value); err != nil {
		badRequest(w, "invalid request body")
		return
	}

	filename, err := thirdparty.UploadFile(ctx, file, filename)
	if err != nil {
		internalServerError(w)
		return
	}

	body.Image = filename

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
	mediatype, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	body := make(map[string]interface{})

	if mediatype == "multipart/form-data" {
		if err := r.ParseMultipartForm(5 << 20); err != nil {
			badRequest(w, "maximum upload size is 5 MB")
			return
		}

		for k, v := range r.MultipartForm.Value {
			body[k] = v[0]
		}

		file, header, fileErr := r.FormFile("image")
		if fileErr != nil {
			badRequest(w, fileErr.Error())
			return
		}
		defer file.Close()

		if err := helper.ValidateImageFile(header.Header.Get("Content-Type")); err != nil {
			badRequest(w, "upload only for image")
			return
		}

		topping, err := s.ProductUseCase.GetTopping(ctx, toppingID)
		if err != nil {
			internalServerError(w)
			return
		}

		filename, err := thirdparty.UpdateImage(file, topping.Image, header.Filename)
		if err != nil {
			internalServerError(w)
			return
		}

		body["image"] = filename

		if err := s.ProductUseCase.UpdateTopping(ctx, toppingID, body); err != nil {
			internalServerError(w)
			return
		}

		resBody, _ := json.Marshal(commonResponse{
			Message: "resource has successfully updated",
		})

		responseOK(w, resBody)
		return
	}

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
	_, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	topping, err := s.ProductUseCase.GetTopping(ctx, toppingID)
	if err != nil {
		notFound(w)
		return
	}

	if err := s.ProductUseCase.DeleteTopping(ctx, toppingID); err != nil {
		internalServerError(w)
		return
	}
	if err := thirdparty.RemoveFile(ctx, topping.Image); err != nil {
		internalServerError(w)
		return
	}

	resBody, _ := json.Marshal(commonResponse{
		Message: "resource has successfully deleted",
	})

	responseOK(w, resBody)
}
