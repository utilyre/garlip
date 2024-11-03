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
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrUsernameTaken   = errors.New("username already taken")
	ErrTokenInvalid    = errors.New("token invalid")
	ErrTokenExpired    = errors.New("token expired")
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

	err = a.Queries.CreateAccount(ctx, queries.CreateAccountParams{
		Username: params.Username,
		Password: hash,
		Fullname: params.Fullname,
	})
	if pgErr := (&pgconn.PgError{}); errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrUsernameTaken
	}
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}

	return nil
}

type AuthLoginParams struct {
	Username string
	Password []byte
}

type JWTClaims struct {
	ID       int32
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

	authInfo, err := a.Queries.GetAccountAuthInfo(ctx, params.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrAccountNotFound
	}
	if err != nil {
		return "", fmt.Errorf("database: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(authInfo.Password, params.Password)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", ErrAccountNotFound
	}
	if err != nil {
		return "", fmt.Errorf("bcrypt: %w", err)
	}

	claims := &JWTClaims{
		ID:       authInfo.ID,
		Username: params.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
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
			return nil, ErrTokenInvalid
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := t.Claims.(*JWTClaims)
	if !ok {
		return "", ErrTokenInvalid
	}
	if time.Now().After(claims.ExpiresAt.Time) {
		return "", ErrTokenExpired
	}

	return claims.Username, nil
}
