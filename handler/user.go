package handler

import (
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

// type UserHandler interface {
// 	GetUsers(w http.ResponseWriter, r *http.Request)
// 	GetUser(w http.ResponseWriter, r *http.Request)
// 	UpdateUser(w http.ResponseWriter, r *http.Request)
// 	DeleteUser(w http.ResponseWriter, r *http.Request)
// 	Register(w http.ResponseWriter, r *http.Request)
// 	Login(w http.ResponseWriter, r *http.Request)
// }

type UserHandler struct {
	UserUseCase usecase.UserUseCase
}

type (
	Register_Req struct {
		Name string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=16"`
		Gender string `json:"gender" validate:"required"`
		Phone string `json:"phone" validate:"required"`
		IsAdmin bool `json:"is_admin"`
	}
	Register_Payload struct {
		Name string `json:"name"`
		Email string `json:"email"`
	}
	Register_Res struct {
		CommonResponse
		Payload Register_Payload `json:"payload"`
	}
)

type (
	Login_Req struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	
	Login_Payload struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Token string `json:"token"`
	}
	Login_Res struct {
		CommonResponse
		Payload Login_Payload `json:"payload"`
	}
)

func (s *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte("get all users"))
}

func (s *UserHandler) GetUser(w http.ResponseWriter, r *http.Request)  {
	userID, _:= strconv.Atoi(chi.URLParam(r, "userID"))
	user, err := s.UserUseCase.FindUserById(r.Context(), userID)

	if err != nil {
		notFound(w)
    return
	}
	responseStruct := struct{
		CommonResponse
		Payload *entity.User `json:"payload"`
	} {
		CommonResponse: CommonResponse{
			Message: "resource has successfully get",
		},
		Payload: user,
	}
	resBody, _ := json.Marshal(responseStruct)
	responseOK(w, resBody)
}

func (s *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _:= strconv.Atoi(chi.URLParam(r, "userID"))

	if claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims); ok {
		if claims.UserID != userID {
			forbidden(w)
			return
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "request invalid")
		return
	}

	if err := s.UserUseCase.UpdateUser(ctx, userID, body); err != nil {
		internalServerError(w)
		return
	}

	resp, _ := json.Marshal(CommonResponse{
		Message: "resource successfully updated",
	})
	responseOK(w, resp)
}

func (s *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := strconv.Atoi(chi.URLParam(r, "userID"))
	
	if claims, ok := ctx.Value(middleware.TokenCtxKey).(*helper.MyClaims); ok {
		if claims.UserID != userID {
			forbidden(w)
			return
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := s.UserUseCase.DeleteUser(ctx, userID); err != nil {
		internalServerError(w)
    return
	}
	
	resBody, _ := json.Marshal(CommonResponse{
		Message:  "resource successfully deleted",
	})
	responseOK(w, resBody)
}

func (s *UserHandler) Register(w http.ResponseWriter, r *http.Request)  {
	ctx := r.Context()
	var body Register_Req
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

	// bytes, encryptErr := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	// if encryptErr != nil {
	// 	internalServerError(w)
  //   return
	// }

	// hashedPassword := string(bytes)
	newUser := entity.User{
		Name: body.Name,
		Email: body.Email,
		Password: body.Password,
		Gender: body.Gender,
		Phone: body.Phone,
		IsAdmin: body.IsAdmin,
	}
	
	if err := s.UserUseCase.CreateNewUser(ctx, newUser); err != nil {
		internalServerError(w)
		return
	}
	responseStruct := Register_Res{
		CommonResponse: CommonResponse{
			Message: "resource successfully created",
		},
		Payload: Register_Payload{
			Name: body.Name,
			Email: body.Email,
		},
	}
	resBody, _ := json.Marshal(responseStruct)
	responseOK(w, resBody)
}

// Handle login from client
// If email not found in DB will return message error with code 404
// If password is not match when compare with hashedPassword in DB
// will return message error with code 400  
func (s *UserHandler) Login(w http.ResponseWriter, r *http.Request)  {
	var body Login_Req
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
		log.Println(err)
		badRequest(w, "credential is not valid")
    return
	}

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, helper.MyClaims{
	// 	UserID: user.Id,
	// 	IsAdmin: user.IsAdmin,
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
	// 		Issuer: "Waysbucks",
	// 	},
	// })

	// secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, tokenErr := helper.GenerateToken(user)
	if tokenErr != nil {
		log.Println(tokenErr)
		return
	}

	responseStruct := Login_Res{
		CommonResponse: CommonResponse{
			Message: "login success",
		},
		Payload: Login_Payload{
			Name: user.Name,
			Email: body.Email,
			Token: tokenString,
		},
	}

	resBody, _ := json.Marshal(responseStruct)
	responseOK(w, resBody)
}