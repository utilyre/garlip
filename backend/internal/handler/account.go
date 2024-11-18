package handler

import (
	"encoding/json"
	"errors"
	"garlip/internal/service"
	"net/http"

	"github.com/utilyre/xmate/v3"
)

type AccountHandler struct {
	AccountSVC *service.AccountService
}

func (ah *AccountHandler) DeleteMe(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return Errorf(http.StatusBadRequest, "Unsupported content type")
	}

	claims := GetClaims(r)
	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return Errorf(http.StatusBadRequest, "Decoding JSON failed due to %v", err)
	}

	err := ah.AccountSVC.DeleteByUsername(r.Context(), service.AccountDeleteParams{
		Username: claims.Username,
		Password: []byte(body.Password),
	})
	if errors.Is(err, service.ErrAccountNotFound) {
		return Errorf(http.StatusNotFound, "Account not found")
	}
	if err != nil {
		return err
	}

	return xmate.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Account has been deleted",
	})
}
