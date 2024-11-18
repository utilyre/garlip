package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"garlip/internal/config"
	"garlip/internal/queries"
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
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("field %s: %s", e.Field, e.Message)
}

type AuthService struct {
	Queries *queries.Queries
}

type AuthRegisterParams struct {
	Username string
	Password []byte
	Fullname string
}

var reUsername = regexp.MustCompile("[0-9A-Za-z_]*")

func (a *AuthService) Register(ctx context.Context, params AuthRegisterParams) error {
	if len(params.Username) == 0 {
		return ValidationError{
			Field:   "username",
			Message: "Required",
		}
	} else if len(params.Username) < 3 {
		return ValidationError{
			Field:   "username",
			Message: "Shorter than 3 chars",
		}
	} else if len(params.Username) > 50 {
		return ValidationError{
			Field:   "username",
			Message: "Longer than 50 chars",
		}
	} else if !reUsername.MatchString(params.Username) {
		return ValidationError{
			Field:   "username",
			Message: "Contains chars other than alphanumeric and underscore",
		}
	}
	if len(params.Password) == 0 {
		return ValidationError{
			Field:   "password",
			Message: "Required",
		}
	} else if len(params.Password) < 8 {
		return ValidationError{
			Field:   "password",
			Message: "Shorter than 8 chars",
		}
	} else if len(params.Password) > 1024 {
		return ValidationError{
			Field:   "password",
			Message: "Longer than 1024 chars",
		}
	}
	if len(params.Fullname) > 100 {
		return ValidationError{
			Field:   "fullname",
			Message: "Longer than 100 chars",
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
	ID       int32  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (a *AuthService) Login(ctx context.Context, params AuthLoginParams) (token string, err error) {
	if len(params.Username) == 0 {
		return "", ValidationError{
			Field:   "username",
			Message: "Required",
		}
	} else if len(params.Username) < 3 {
		return "", ValidationError{
			Field:   "username",
			Message: "Shorter than 3 chars",
		}
	} else if len(params.Username) > 50 {
		return "", ValidationError{
			Field:   "username",
			Message: "Longer than 50 chars",
		}
	} else if !reUsername.MatchString(params.Username) {
		return "", ValidationError{
			Field:   "username",
			Message: "Contains chars other than alphanumeric and underscore",
		}
	}
	if len(params.Password) == 0 {
		return "", ValidationError{
			Field:   "password",
			Message: "Required",
		}
	} else if len(params.Password) < 8 {
		return "", ValidationError{
			Field:   "password",
			Message: "Shorter than 8 chars",
		}
	} else if len(params.Password) > 1024 {
		return "", ValidationError{
			Field:   "password",
			Message: "Longer than 1024 chars",
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Default().TokenLifespan)),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(config.Default().TokenSecret)
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}

	return token, nil
}

type Claims struct {
	ID       int32
	Username string
}

func (a *AuthService) VerifyToken(ctx context.Context, token string) (*Claims, error) {
	var claims JWTClaims
	t, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok ||
			method != jwt.SigningMethodHS256 {
			return nil, ErrTokenInvalid
		}

		return config.Default().TokenSecret, nil
	})
	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, ErrTokenExpired
	}
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, ErrTokenInvalid
	}

	return &Claims{
		ID:       claims.ID,
		Username: claims.Username,
	}, nil
}
