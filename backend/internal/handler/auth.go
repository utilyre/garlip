package handler

import (
	"encoding/json"
	"errors"
	"garlip/internal/service"
	"net/http"

	"github.com/utilyre/xmate/v2"
)

type AuthHandler struct {
	AuthSVC *service.AuthService
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return xmate.Errorf(http.StatusBadRequest, "Unsupported content type")
	}

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Fullname string `json:"fullname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return xmate.Errorf(http.StatusBadRequest, "Decoding JSON failed due to %v", err)
	}

	err := a.AuthSVC.Register(r.Context(), service.AuthRegisterParams{
		Username: body.Username,
		Password: []byte(body.Password),
		Fullname: body.Fullname,
	})
	if errors.Is(err, service.ErrAccountDup) {
		return xmate.Errorf(http.StatusConflict, "Account already exists")
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
		return xmate.Errorf(http.StatusBadRequest, "Unsupported content type")
	}

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return xmate.Errorf(http.StatusBadRequest, "Decoding JSON failed due to %v", err)
	}

	token, err := a.AuthSVC.Login(r.Context(), service.AuthLoginParams{
		Username: "",
		Password: []byte{},
	})
	if err != nil {
		return err
	}

	return xmate.WriteJSON(w, http.StatusOK, map[string]any{
		"token": token,
	})
}
