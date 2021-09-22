package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/yosepalexsander/waysbucks-api/domain"
	"github.com/yosepalexsander/waysbucks-api/persistance"

	"golang.org/x/crypto/bcrypt"
)

type UserServer struct {
	Finder persistance.UserFinder
	Saver persistance.UserSaver
	Remover persistance.UserRemover
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
	userID, _:= strconv.ParseUint(chi.URLParam(r, "userID"), 10, 32)
	user, err := s.Finder.FindUserById(r.Context(), userID)

	if err != nil {
		notFound(w)
    return
	}
	responseStruct := struct{
		CommonResponse
		Payload *domain.User `json:"payload"`
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
	userID, _:= strconv.ParseUint(chi.URLParam(r, "userID"), 10, 32)
	claims, ok := ctx.Value(tokenCtxKey).(*MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	
	if claims.UserID != userID {
		forbidden(w)
		return
	}

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "request invalid")
		return
	}
	
	_, err := s.Saver.UpdateUser(ctx, claims.UserID, body)

	if err != nil {
		internalServerError(w)
		return
	}
}

func (s *UserServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 32)
	claims, ok := ctx.Value(tokenCtxKey).(*MyClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if claims.UserID != userID {
		forbidden(w)
		return
	}

	if err := s.Remover.DeleteUser(ctx, userID); err != nil {
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
		w.WriteHeader(http.StatusBadRequest)
    return
	}

	if user, _ := s.Finder.FindUserByEmail(r.Context(), body.Email); user.Email == body.Email {
		badRequest(w, "resource already exist")
    return
	}

	isValid, msg := validate(body)
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
	newUser := domain.User{
		Name: body.Name,
		Email: body.Email,
		Password: hashedPassword,
		Gender: body.Gender,
		Phone: body.Phone,
		IsAdmin: body.IsAdmin,
	}
	
	if err := s.Saver.SaveUser(r.Context(), newUser); err != nil {
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

	user, findErr := s.Finder.FindUserByEmail(r.Context(), body.Email)

	if findErr != nil {
		notFound(w)
    return
	}

	hashedPassword := []byte(user.Password)
	reqPassword := []byte(body.Password)
	if err := bcrypt.CompareHashAndPassword(hashedPassword, reqPassword); err != nil {
		badRequest(w, "credential is not valid")
    return
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims{
		user.Id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
			Issuer: "Waysbucks",
		},
	})

	// Sign and get the complete encoded token as a string using the secret
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


func validate(value interface{}) (bool, string)  {
	v := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	addTranslation(v, trans, "email", "{0} must be a valid email")
	addTranslation(v, trans, "min", "{0} must be at least {1} char length")
	addTranslation(v, trans, "max", "{0} must be max {1} char length")
	addTranslation(v, trans, "required", "{0} is a required field")
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	err := v.Struct(value)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, err := range validationErrors {
			log.Println(err.Error())
		}
		msgErr := validationErrors[0].Translate(trans)
		return false, msgErr
	}
	return true, ""
}

func addTranslation(v *validator.Validate, trans ut.Translator, tag string, errMessage string) {
	registerFn := func(ut ut.Translator) error {
		 return ut.Add(tag, errMessage, false)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		 tag := fe.Tag()

		 t, err := ut.T(tag, fe.Field(), param)
		 if err != nil {
				return fe.(error).Error()
		 }
		 return t
	}
	_ = v.RegisterTranslation(tag, trans, registerFn, transFn)
}