package handler

import (
	"encoding/json"
	"garlip/internal/service"
	"log"
	"net/http"
)

type AuthHandler struct {
	AuthSVC *service.AuthService
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Fullname string `json:"fullname"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := a.AuthSVC.Register(r.Context(), service.AuthRegisterParams{
		Username: req.Username,
		Password: []byte(req.Password),
		Fullname: req.Fullname,
	}); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"message": "account successfully registered",
	}); err != nil {
		log.Println(err)
	}
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	token, err := a.AuthSVC.Login(r.Context(), service.AuthLoginParams{
		Username: "",
		Password: []byte{},
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"token": token,
	}); err != nil {
		log.Println(err)
	}
}
