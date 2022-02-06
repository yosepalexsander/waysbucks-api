package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/handler/middleware"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
	"github.com/yosepalexsander/waysbucks-api/usecase"
)

type UserHandler struct {
	UserUseCase usecase.UserUseCase
}

func NewUserHandler(u usecase.UserUseCase) UserHandler {
	return UserHandler{u}
}

func (s *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	type response struct {
		commonResponse
		Payload []entity.User `json:"payload"`
	}

	users, err := s.UserUseCase.FindUsers(r.Context())

	if err != nil {
		internalServerError(w)
	}

	responseStruct := response{
		commonResponse: commonResponse{
			Message: "get resources successfully",
		},
		Payload: users,
	}

	resp, _ := json.Marshal(responseStruct)
	responseOK(w, resp)
}

func (s *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	type response struct {
		commonResponse
		Payload *entity.User `json:"payload"`
	}

	if claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims); ok {
		user, err := s.UserUseCase.GetProfile(ctx, claims.UserID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				notFound(w)
			default:
				internalServerError(w)
			}
			return
		}

		if user.Image, err = thirdparty.GetImageUrl(ctx, user.Image); err != nil {
			serviceUnavailable(w, "error: cloudinary service unavailable")
			return
		}

		responseStruct := response{
			commonResponse: commonResponse{
				Message: "resource has successfully get",
			},
			Payload: user,
		}
		resBody, _ := json.Marshal(responseStruct)

		responseOK(w, resBody)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func (s *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims); ok {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			badRequest(w, "request invalid")
			return
		}

		if err := s.UserUseCase.UpdateUser(ctx, claims.UserID, body); err != nil {
			internalServerError(w)
			return
		}

		resp, _ := json.Marshal(commonResponse{
			Message: "resource successfully updated",
		})

		responseOK(w, resp)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func (s *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims); ok {
		body := make(map[string]interface{})

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

		user, err := s.UserUseCase.GetProfile(ctx, claims.UserID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				notFound(w)
			default:
				internalServerError(w)
			}

			return
		}

		filename, err := thirdparty.UpdateImage(file, user.Image, header.Filename)
		if err != nil {
			internalServerError(w)
			return
		}

		body["image"] = filename
		if err := s.UserUseCase.UpdateUser(ctx, claims.UserID, body); err != nil {
			internalServerError(w)
			return
		}

		resp, _ := json.Marshal(commonResponse{
			Message: "resource has successfully created",
		})

		responseOK(w, resp)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func (s *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims); ok {
		if err := s.UserUseCase.DeleteUser(ctx, claims.UserID); err != nil {
			internalServerError(w)
			return
		}

		resBody, _ := json.Marshal(commonResponse{
			Message: "resource successfully deleted",
		})

		responseOK(w, resBody)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

}

func (s *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	type (
		request struct {
			Name     string `json:"name" validate:"required"`
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required,min=8,max=16"`
			Gender   string `json:"gender" validate:"required"`
			Phone    string `json:"phone" validate:"required"`
			IsAdmin  bool   `json:"is_admin"`
		}
		payload struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		response struct {
			commonResponse
			Payload payload `json:"payload"`
		}
	)

	ctx := r.Context()

	var body request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}

	isValid, msg := helper.Validate(body)
	if !isValid {
		badRequest(w, msg)
		return
	}

	if user, _ := s.UserUseCase.FindUserByEmail(ctx, body.Email); user.Email == body.Email {
		badRequest(w, "resource already exist")
		return
	}

	newUser := entity.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
		Gender:   body.Gender,
		Phone:    body.Phone,
		IsAdmin:  body.IsAdmin,
	}

	if err := s.UserUseCase.CreateNewUser(ctx, newUser); err != nil {
		internalServerError(w)
		return
	}

	responseStruct := response{
		commonResponse: commonResponse{
			Message: "resource successfully created",
		},
		Payload: payload{
			Name:  body.Name,
			Email: body.Email,
		},
	}
	resBody, _ := json.Marshal(responseStruct)

	responseOK(w, resBody)
}

func (s *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	type (
		request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		payload struct {
			Id    int    `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
			Token string `json:"token"`
		}
		response struct {
			commonResponse
			Payload payload `json:"payload"`
		}
	)

	var body request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}

	user, err := s.UserUseCase.FindUserByEmail(r.Context(), body.Email)
	if err != nil {
		notFound(w)
		return
	}

	if err := s.UserUseCase.ValidatePassword(user.Password, body.Password); err != nil {
		badRequest(w, "credential is not valid")
		return
	}

	tokenString, tokenErr := helper.GenerateToken(user)
	if tokenErr != nil {
		internalServerError(w)
		return
	}

	responseStruct := response{
		commonResponse: commonResponse{
			Message: "login success",
		},
		Payload: payload{
			Id:    user.Id,
			Name:  user.Name,
			Email: body.Email,
			Token: tokenString,
		},
	}
	resBody, _ := json.Marshal(responseStruct)

	responseOK(w, resBody)
}
