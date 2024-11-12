package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"garlip/internal/service"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return xmate.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "You have been logged in",
	})
}

type Auth struct {
	ID       int32
	Username string
}

func (ah *AuthHandler) Authenticate(next http.Handler) http.Handler {
	return xmate.HandleFunc(func(w http.ResponseWriter, r *http.Request) error {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			return err
		}

		var claims service.JWTClaims
		token, err := jwt.ParseWithClaims(cookie.Value, &claims, func(t *jwt.Token) (any, error) {
			method, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok || method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			return err
		}
		if !token.Valid {
			return errors.New("invalid token")
		}

		r2 := r.WithContext(context.WithValue(r.Context(), "auth", Auth{
			ID:       claims.ID,
			Username: claims.Username,
		}))

		next.ServeHTTP(w, r2)
		return nil
	})
}
