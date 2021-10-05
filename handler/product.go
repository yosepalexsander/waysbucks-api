package handler

import (
	"encoding/json"
	"mime"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/handler/middleware"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/usecase"
)

type ProductHandler struct {
	ProductUseCase usecase.ProductUseCase
}

func (s *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// specifies a maximum upload of 10MB file size
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		badRequest(w, "maximum upload size is 10 MB")
		return
	}

	file, header, fileErr := r.FormFile("image")

	if fileErr != nil {
		badRequest(w, fileErr.Error())
		return 
	}
	
	defer file.Close()

	filename, err := helper.UploadFile(ctx, file, header.Filename)

	if err != nil  {
		if err.Error() == "invalid file extension" {
			badRequest(w, "upload only for image")
			return
		}
		internalServerError(w)
		return
	}
	
	var body entity.Product
 
	if err := schema.NewDecoder().Decode(&body, r.MultipartForm.Value); err != nil {
		badRequest(w, "invalid request")
		return
	}

	body.Image = filename

	if isValid, msg := helper.Validate(body); !isValid {
		badRequest(w, msg)
		return
	}
	if err := s.ProductUseCase.CreateProduct(ctx, body); err != nil {
		internalServerError(w)
		return
	}

	resp, _ := json.Marshal(commonResponse{
		Message: "resource successfully created",
	})
	responseOK(w, resp)
}

func (s *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := s.ProductUseCase.GetAllProduct(r.Context())

	if err != nil {
		internalServerError(w)
		return
	}

	responseStruct := struct{
		commonResponse
		Payload []entity.Product `json:"payload"`
	} {
		commonResponse: commonResponse{
			Message: "resource has successfully get",
		},
		Payload: products,
	}
	resp, _ := json.Marshal(responseStruct)

	responseOK(w, resp)
}

func (s *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID, _ := strconv.Atoi(chi.URLParam(r, "productID"))
	product, err := s.ProductUseCase.GetProduct(ctx, productID)

	if err != nil {
		internalServerError(w)
		return
	}

	responseStruct := struct{
		commonResponse
		Payload *entity.Product `json:"payload"`
	} {
		commonResponse: commonResponse{
			Message: "resource has successfully get",
		},
		Payload: product,
	}

	resp, _ := json.Marshal(responseStruct)
	responseOK(w, resp)
}


func (s *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID, _ := strconv.Atoi(chi.URLParam(r, "productID"))
	mediatype, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))

	body := make(map[string]interface{})
	
	if mediatype == "multipart/form-data" {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			badRequest(w, "maximum upload size is 10 MB")
			return
		}
		
		for k, v := range r.PostForm {
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
			internalServerError(w)
			return
		}

		filename, err := helper.UploadFile(ctx, file, header.Filename)

		if err != nil  {
			if err.Error() == "invalid file extension" {
				badRequest(w, "upload only for image")
				return
			}
			internalServerError(w)
			return
		}
		
		body["image"] = filename

		if err := helper.RemoveFile(ctx, product.Image); err != nil {
			internalServerError(w)
			return
		}

		if err := s.ProductUseCase.UpdateProduct(ctx, productID, body); err != nil {
			internalServerError(w)
			return
		}

		resBody, _ := json.Marshal(commonResponse{
			Message:  "resource successfully deleted",
		})
		responseOK(w, resBody)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}

	if err := s.ProductUseCase.UpdateProduct(ctx, productID, body); err != nil {
		internalServerError(w)
		return
	}

	resBody, _ := json.Marshal(commonResponse{
		Message:  "resource successfully deleted",
	})
	responseOK(w, resBody)
}

func (s *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID, _ := strconv.Atoi(chi.URLParam(r, "productID"))
	
	if _, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims); ok {
		product, err := s.ProductUseCase.GetProduct(ctx, productID)
		if err != nil {
			internalServerError(w)
			return
		}

		if err := helper.RemoveFile(ctx, product.Image); err != nil {
			internalServerError(w)
			return
		}

		if err := s.ProductUseCase.DeleteProduct(ctx, productID); err != nil {
			if err.Error() == "no rows affected" {
				notFound(w)
				return
			}
			internalServerError(w)
			return
		}
		
		resBody, _ := json.Marshal(commonResponse{
			Message:  "resource successfully deleted",
		})
		responseOK(w, resBody)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
}