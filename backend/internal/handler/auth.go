package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"garlip/internal/config"
	"garlip/internal/service"
	"net/http"
	"time"

	"github.com/utilyre/xmate/v3"
)

type Error struct {
	Status  int
	Message string
}

func Errorf(status int, format string, a ...any) Error {
	return Error{
		Status:  status,
		Message: fmt.Sprintf(format, a...),
	}
}

func (he Error) Error() string {
	return he.Message
}

type AuthHandler struct {
	AuthSVC *service.AuthService
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return Errorf(http.StatusBadRequest, "Unsupported content type")
	}

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Fullname string `json:"fullname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return Errorf(http.StatusBadRequest, "Decoding JSON failed due to %v", err)
	}

	err := a.AuthSVC.Register(r.Context(), service.AuthRegisterParams{
		Username: body.Username,
		Password: []byte(body.Password),
		Fullname: body.Fullname,
	})
	if errors.Is(err, service.ErrUsernameTaken) {
		return Errorf(http.StatusConflict, "Account already exists")
	}
	if err != nil {
		return err
	}

	return xmate.WriteJSON(w, http.StatusCreated, map[string]any{
		"message": "Account has been registered",
	})
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return Errorf(http.StatusBadRequest, "Unsupported content type")
	}

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return Errorf(http.StatusBadRequest, "Decoding JSON failed due to %v", err)
	}

	token, err := a.AuthSVC.Login(r.Context(), service.AuthLoginParams{
		Username: body.Username,
		Password: []byte(body.Password),
	})
	if errors.Is(err, service.ErrAccountNotFound) {
		return Errorf(http.StatusNotFound, "Account not found")
	}
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieJWT,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(config.Default().TokenLifespan),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return xmate.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "You have been logged in",
	})
}

const cookieJWT = "jwt"

type Key int

const (
	KeyClaims = iota + 1
)

func (ah *AuthHandler) Authenticate(next http.Handler) http.Handler {
	return xmate.HandleFunc(func(w http.ResponseWriter, r *http.Request) error {
		cookie, err := r.Cookie(cookieJWT)
		if errors.Is(err, http.ErrNoCookie) {
			return Errorf(http.StatusUnauthorized, "Cookie not provided")
		}
		if err != nil {
			return err
		}

		claims, err := ah.AuthSVC.VerifyToken(r.Context(), cookie.Value)
		if errors.Is(err, service.ErrTokenInvalid) {
			return Errorf(http.StatusUnauthorized, "Invalid token")
		}
		if errors.Is(err, service.ErrTokenExpired) {
			return Errorf(http.StatusUnauthorized, "Expired token")
		}
		if err != nil {
			return err
		}

		r2 := r.WithContext(context.WithValue(r.Context(), KeyClaims, claims))
		next.ServeHTTP(w, r2)
		return nil
	})
}
