package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"garlip/internal/queries"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrInvalidToken    = errors.New("invalid token")
	ErrExpiredToken    = errors.New("expired token")
)

type ValidationError struct {
	Param string
	Msg   string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("param %s: %s", e.Param, e.Msg)
}

type AuthService struct {
	Queries *queries.Queries
}

type AuthRegisterParams struct {
	Username string
	Password []byte
	Fullname string
}

func (a *AuthService) Register(ctx context.Context, params AuthRegisterParams) error {
	if len(params.Username) < 3 {
		return ValidationError{
			Param: "username",
			Msg:   "shorter than 3 chars",
		}
	}
	if len(params.Username) > 50 {
		return ValidationError{
			Param: "username",
			Msg:   "longer than 50 chars",
		}
	}
	re := regexp.MustCompile("[0-9A-Za-z_]*")
	if !re.MatchString(params.Username) {
		return ValidationError{
			Param: "username",
			Msg:   "contains chars other than alphanumeric and underscore",
		}
	}
	if len(params.Password) < 8 {
		return ValidationError{
			Param: "password",
			Msg:   "shorter than 8 chars",
		}
	}
	if len(params.Password) > 1024 {
		return ValidationError{
			Param: "password",
			Msg:   "longer than 1024 chars",
		}
	}
	if len(params.Fullname) > 100 {
		return ValidationError{
			Param: "fullname",
			Msg:   "longer than 100 chars",
		}
	}

	hash, err := bcrypt.GenerateFromPassword(params.Password, bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt: %w", err)
	}

	if err := a.Queries.CreateAccount(ctx, queries.CreateAccountParams{
		Username: params.Username,
		Password: hash,
		Fullname: params.Fullname,
	}); err != nil {
		// TODO: handle dup case
		return fmt.Errorf("database: %w", err)
	}

	return nil
}

type AuthLoginParams struct {
	Username string
	Password []byte
}

type JWTClaims struct {
	Username string
	jwt.RegisteredClaims
}

func (a *AuthService) Login(ctx context.Context, params AuthLoginParams) (token string, err error) {
	if len(params.Username) < 3 {
		return "", ValidationError{
			Param: "username",
			Msg:   "shorter than 3 chars",
		}
	}
	if len(params.Username) > 50 {
		return "", ValidationError{
			Param: "username",
			Msg:   "longer than 50 chars",
		}
	}
	re := regexp.MustCompile("[0-9A-Za-z_]*")
	if !re.MatchString(params.Username) {
		return "", ValidationError{
			Param: "username",
			Msg:   "contains chars other than alphanumeric and underscore",
		}
	}
	if len(params.Password) < 8 {
		return "", ValidationError{
			Param: "password",
			Msg:   "shorter than 8 chars",
		}
	}
	if len(params.Password) > 1024 {
		return "", ValidationError{
			Param: "password",
			Msg:   "longer than 1024 chars",
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

	claims := &JWTClaims{
		Username: params.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodES256, claims).
		SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}

	return token, nil
}

func (a *AuthService) VerifyToken(ctx context.Context, token string) (username string, err error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if method, ok := t.Method.(*jwt.SigningMethodECDSA); !ok ||
			method != jwt.SigningMethodES256 {
			return nil, ErrInvalidToken
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := t.Claims.(*JWTClaims)
	if !ok {
		return "", ErrInvalidToken
	}
	if time.Now().After(claims.ExpiresAt.Time) {
		return "", ErrExpiredToken
	}

	return claims.Username, nil
}
