package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"garlip/internal/postgres"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

type ValidationError struct {
	Param string
	Err   error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("param %s: %v", e.Param, e.Err)
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

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
		return ValidationError{
			Param: "username",
			Err:   errors.New("shorter than 3 chars"),
		}
	}
	if len(params.Username) > 50 {
		return ValidationError{
			Param: "username",
			Err:   errors.New("longer than 50 chars"),
		}
	}
	re := regexp.MustCompile("[0-9A-Za-z_]*")
	if !re.MatchString(params.Username) {
		return ValidationError{
			Param: "username",
			Err:   errors.New("contains chars other than alphanumeric and underscore"),
		}
	}

	if len(params.Password) < 8 {
		return ValidationError{
			Param: "password",
			Err:   errors.New("short than 8 chars"),
		}
	}

	hash, err := bcrypt.GenerateFromPassword(params.Password, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt: %w", err)
	}

	if err := a.Queries.CreateAccount(ctx, postgres.CreateAccountParams{
		Username: params.Username,
		Password: hash,
		Fullname: sql.NullString{String: params.Fullname, Valid: true},
	}); err != nil {
		return fmt.Errorf("database: %w", err)
	}

	return nil
}

type AuthLoginParams struct {
	Username string
	Password []byte
}

type JWTClaims struct {
	jwt.RegisteredClaims
	Username string
}

func (a *Auth) Login(ctx context.Context, params AuthLoginParams) (string, error) {
	if len(params.Username) < 3 {
		return "", ValidationError{
			Param: "username",
			Err:   errors.New("shorter than 3 chars"),
		}
	}
	if len(params.Username) > 50 {
		return "", ValidationError{
			Param: "username",
			Err:   errors.New("longer than 50 chars"),
		}
	}
	re := regexp.MustCompile("[0-9A-Za-z_]*")
	if !re.MatchString(params.Username) {
		return "", ValidationError{
			Param: "username",
			Err:   errors.New("contains chars other than alphanumeric and underscore"),
		}
	}

	if len(params.Password) < 8 {
		return "", ValidationError{
			Param: "password",
			Err:   errors.New("short than 8 chars"),
		}
	}

	hash, err := a.Queries.GetAccountPasswordByUsername(ctx, params.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrAccountNotFound
	}
	if err != nil {
		return "", fmt.Errorf("database: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(hash, params.Password)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", ErrAccountNotFound
	}
	if err != nil {
		return "", fmt.Errorf("bcrypt: %w", err)
	}

	var claims JWTClaims
	claims.Username = params.Username
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))

	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).
		SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}

	return token, nil
}
