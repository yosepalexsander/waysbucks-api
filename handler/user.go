package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/yosepalexsander/waysbucks-api/domain"
	"github.com/yosepalexsander/waysbucks-api/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserServer struct {
	Finder storage.UserFinder
	Saver storage.UserSaver
	Delete storage.UserDelete
}

type CommonResponse struct {
	Message string `json:"message"`
}

type (
	Register_Req struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Password string `json:"password"`
		Gender string `json:"gender"`
		Phone string `json:"phone"`
		IsAdmin uint8
	}
	Register_Payload struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Token string `json:"token"`
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

func (s UserServer) GetUsers(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte("get all users"))
}

func (s UserServer) GetUser(w http.ResponseWriter, r *http.Request)  {
	userID, _:= strconv.ParseUint(chi.URLParam(r, "userID"), 10, 32)
	user, err := s.Finder.FindUserById(r.Context(), userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(CommonResponse{
			Message:  "error",
		})
		w.Write(resp)
    return
	}
	res_body := struct{
		CommonResponse
		Payload *domain.User `json:"payload"`
	} {
		CommonResponse: CommonResponse{
			Message: "resource has successfully get",
		},
		Payload: user,
	}
	resp, _ := json.Marshal(res_body)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (s UserServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 32)

	if err := s.Delete.DeleteUser(r.Context(), userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(CommonResponse{
			Message:  "error",
		})
		w.Write(resp)
    return
	}
	
	resp, _ := json.Marshal(CommonResponse{
		Message:  "resource successfully deleted",
	})
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (s UserServer) Register(w http.ResponseWriter, r *http.Request)  {
	var body Register_Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
    return
	}

	if user, _ := s.Finder.FindUserByEmail(r.Context(), body.Email); user.Email == body.Email {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(CommonResponse{
			Message: "resource already exist",
		})
		w.Write(resp)
    return
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(CommonResponse{
			Message: "error",
		})
		w.Write(resp)
    return
	}
	hashedPassword := string(bytes)
	newUser := domain.User{
		Name: body.Name,
		Email: body.Email,
		Password: hashedPassword,
		Gender: body.Gender,
		Phone: body.Phone,
		IsAdmin: body.IsAdmin,
	}
	
	if err := s.Saver.SaveUser(r.Context(), newUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("database error"))
		return
	}
	response := Register_Res{
		CommonResponse: CommonResponse{
			Message: "resource successfully created",
		},
		Payload: Register_Payload{
			Name: body.Name,
			Email: body.Email,
		},
	}
	resp, _ := json.Marshal(response)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// Handle login from client
// If email not found in DB will return message error with code 404
// If password is not match when compare with hashedPassword in DB
// will return message error with code 400  
func (s UserServer) Login(w http.ResponseWriter, r *http.Request)  {
	var body Login_Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
    return
	}

	user, findErr := s.Finder.FindUserByEmail(r.Context(), body.Email)

	if findErr != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("error"))
    return
	}
	hashedPassword := []byte(user.Password)
	reqPassword := []byte(body.Password)
	if err := bcrypt.CompareHashAndPassword(hashedPassword, reqPassword); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("credential is not valid"))
    return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"owner_id": user.Id,
		"exp": 18000,
		"iss": "user auth",
	})

	// Sign and get the complete encoded token as a string using the secret
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, tokenErr := token.SignedString(secretKey)
	if tokenErr != nil {
		log.Println(tokenErr)
	}

	response := Login_Res{
		CommonResponse: CommonResponse{
			Message: "resource successfully created",
		},
		Payload: Login_Payload{
			Name: user.Name,
			Email: body.Email,
			Token: tokenString,
		},
	}

	resp, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}