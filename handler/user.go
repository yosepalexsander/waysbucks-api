package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yosepalexsander/waysbucks-api/config"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/middleware"
	"github.com/yosepalexsander/waysbucks-api/usecase"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
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

	if user, err := s.UserUseCase.GetUserByEmail(ctx, body.Email); err == nil && user.Email == body.Email {
		badRequest(w, "resource already exist")
		return
	}

	user, err := s.UserUseCase.CreateNewUser(ctx, body.Name, body.Email, body.Password, body.Gender, body.Phone)
	if err != nil {
		internalServerError(w)
		return
	}

	tokenString, tokenErr := helper.GenerateToken(user.Id, user.IsAdmin)
	if tokenErr != nil {
		internalServerError(w)
		return
	}

	responseStruct := AuthResponse{
		commonResponse: commonResponse{
			Message: "resource successfully created",
		},
		Payload: AuthResponsePayload{
			Name:  user.Name,
			Email: user.Email,
			Token: tokenString,
		},
	}
	resBody, _ := json.Marshal(responseStruct)

	responseOK(w, resBody)
}

func (s *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	type (
		request struct {
			Email    string `json:"email" validate:"required"`
			Password string `json:"password" validate:"required"`
		}
	)

	var body request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid request")
		return
	}

	user, err := s.UserUseCase.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		notFound(w)
		return
	}

	isValid, msg := helper.Validate(body)
	if !isValid {
		badRequest(w, msg)
		return
	}

	if err := s.UserUseCase.ValidatePassword(user.Password, body.Password); err != nil {
		badRequest(w, err.Error())
		return
	}

	tokenString, tokenErr := helper.GenerateToken(user.Id, user.IsAdmin)
	if tokenErr != nil {
		internalServerError(w)
		return
	}

	responseStruct := AuthResponse{
		commonResponse: commonResponse{
			Message: "login success",
		},
		Payload: AuthResponsePayload{
			Name:  user.Name,
			Email: user.Email,
			Token: tokenString,
		},
	}
	resBody, _ := json.Marshal(responseStruct)

	responseOK(w, resBody)
}

func (s *UserHandler) LoginOrRegisterWithGoogle(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Token string `json:"token" validate:"required"`
	}

	ctx := r.Context()

	var body request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "Invalid request")
		return
	}

	isValid, msg := helper.Validate(body)
	if !isValid {
		badRequest(w, msg)
		return
	}

	userInfo, err := VerifyTokenID(body.Token)
	if err != nil {
		badRequest(w, "Invalid request")
		return
	}

	if userInfo.Aud != config.GOOGLE_CLIENT_ID {
		badRequest(w, "Google client is not valid")
		return
	}

	user, err := s.UserUseCase.GetUserByEmail(ctx, userInfo.Email)
	if err != nil && err != sql.ErrNoRows {
		internalServerError(w)
		return
	}

	if user == nil {
		user, err = s.UserUseCase.CreateNewUser(ctx, userInfo.Name, userInfo.Email, "", "", "")
		if err != nil {
			internalServerError(w)
			return
		}

	}

	tokenString, tokenErr := helper.GenerateToken(user.Id, user.IsAdmin)
	if tokenErr != nil {
		internalServerError(w)
		return
	}

	responseStruct := AuthResponse{
		commonResponse: commonResponse{
			Message: "login or register with google success",
		},
		Payload: AuthResponsePayload{
			Name:  user.Name,
			Email: user.Email,
			Token: tokenString,
		},
	}
	resBody, _ := json.Marshal(responseStruct)

	responseOK(w, resBody)
}

func VerifyTokenID(idToken string) (*TokenInfo, error) {
	authService, err := oauth2.NewService(context.TODO(), option.WithHTTPClient(http.DefaultClient))
	if err != nil {
		return nil, err
	}

	tokenInfoCall := authService.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancelFunc()
	tokenInfoCall.Context(ctx)

	_, err = tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}

	token, _, err := new(jwt.Parser).ParseUnverified(idToken, &TokenInfo{})
	if tokenInfo, ok := token.Claims.(*TokenInfo); ok {
		tokenInfo.ExpiresAt = tokenInfo.Exp
		if err := tokenInfo.Valid(); err != nil {
			return nil, err
		}

		return tokenInfo, nil
	}

	return nil, err
}

type (
	AuthResponse struct {
		commonResponse
		Payload AuthResponsePayload `json:"payload"`
	}
	AuthResponsePayload struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Token string `json:"token"`
	}
)

type TokenInfo struct {
	Iss string `json:"iss"`
	// userId
	Sub string `json:"sub"`
	Azp string `json:"azp"`
	// clientId
	Aud string `json:"aud"`
	Iat int64  `json:"iat"`
	// expired time
	Exp int64 `json:"exp"`

	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	jwt.StandardClaims
}
