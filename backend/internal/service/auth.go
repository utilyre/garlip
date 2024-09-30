package service

import (
	"context"
	"database/sql"
	"errors"
	"garlip/internal/postgres"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Queries *postgres.Queries
}

type AuthRegisterParams struct {
	Username string
	Password []byte
	Fullname string
}

func (a *Auth) Register(ctx context.Context, params AuthRegisterParams) error {
	if len(params.Username) < 3 {
		return errors.New("username: too short")
	}
	if len(params.Username) > 50 {
		return errors.New("username: too long")
	}
	re := regexp.MustCompile("[0-9A-Za-z]*")
	if !re.MatchString(params.Username) {
		return errors.New("username: regex mismatch")
	}

	if len(params.Password) < 8 {
		return errors.New("password: too short")
	}

	hash, err := bcrypt.GenerateFromPassword(params.Password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return a.Queries.CreateAccount(ctx, postgres.CreateAccountParams{
		Username: params.Username,
		Password: hash,
		Fullname: sql.NullString{String: params.Fullname, Valid: true},
	})
}

type AuthLoginParams struct {
	Username string
	Password []byte
}

func (a *Auth) Login(ctx context.Context, params AuthLoginParams) (string, error) {
	return "", nil
}
