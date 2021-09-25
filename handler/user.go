package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/handler/middleware"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/persistance"

	"golang.org/x/crypto/bcrypt"
)

type UserServer struct {
	Repo persistance.UserRepository
}

type CommonResponse struct {
	Message string `json:"message"`
}

type (
	Register_Req struct {
		Name string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=16"`
		Gender string `json:"gender" validate:"required"`
		Phone string `json:"phone" validate:"required"`
		IsAdmin uint8
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

func (s *UserServer) GetUsers(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte("get all users"))
}

func (s *UserServer) GetUser(w http.ResponseWriter, r *http.Request)  {
	userID, _:= strconv.Atoi(chi.URLParam(r, "userID"))
	user, err := s.Repo.FindUserById(r.Context(), userID)

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

func (s *UserServer) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _:= strconv.Atoi(chi.URLParam(r, "userID"))

	if claims, ok := ctx.Value(middleware.TokenCtxKey).(*middleware.MyClaims); ok {
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

	if err := s.Repo.UpdateUser(ctx, userID, body); err != nil {
		internalServerError(w)
		return
	}

	resp, _ := json.Marshal(CommonResponse{
		Message: "resource successfully updated",
	})
	responseOK(w, resp)
}

func (s *UserServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := strconv.Atoi(chi.URLParam(r, "userID"))
	
	if claims, ok := ctx.Value(middleware.TokenCtxKey).(*middleware.MyClaims); ok {
		if claims.UserID != userID {
			forbidden(w)
			return
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := s.Repo.DeleteUser(ctx, userID); err != nil {
		internalServerError(w)
    return
	}
	
	resBody, _ := json.Marshal(CommonResponse{
		Message:  "resource successfully deleted",
	})
	responseOK(w, resBody)
}

func (s *UserServer) Register(w http.ResponseWriter, r *http.Request)  {
	var body Register_Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
    return
	}

	if user, _ := s.Repo.FindUserByEmail(r.Context(), body.Email); user.Email == body.Email {
		badRequest(w, "resource already exist")
    return
	}

	isValid, msg := helper.Validate(body)
	if !isValid {
		badRequest(w, msg)
		return
	}

	bytes, encryptErr := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if encryptErr != nil {
		internalServerError(w)
    return
	}

	hashedPassword := string(bytes)
	newUser := entity.User{
		Name: body.Name,
		Email: body.Email,
		Password: hashedPassword,
		Gender: body.Gender,
		Phone: body.Phone,
		IsAdmin: body.IsAdmin,
	}
	
	if err := s.Repo.SaveUser(r.Context(), newUser); err != nil {
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
func (s *UserServer) Login(w http.ResponseWriter, r *http.Request)  {
	var body Login_Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "request is not valid")
    return
	}

	user, Repor := s.Repo.FindUserByEmail(r.Context(), body.Email)

	if Repor != nil {
		notFound(w)
    return
	}

	hashedPassword := []byte(user.Password)
	reqPassword := []byte(body.Password)
	if err := bcrypt.CompareHashAndPassword(hashedPassword, reqPassword); err != nil {
		badRequest(w, "credential is not valid")
    return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, middleware.MyClaims{
		UserID: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
			Issuer: "Waysbucks",
		},
	})

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, tokenErr := token.SignedString(secretKey)
	if tokenErr != nil {
		log.Println(tokenErr)
	}

	responseStruct := Login_Res{
		CommonResponse: CommonResponse{
			Message: "resource successfully created",
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