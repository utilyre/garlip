package service

import (
	"context"
	"database/sql"
	"errors"
	"garlip/internal/queries"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	Queries *queries.Queries
}

type AccountUpdateByIDParams struct {
	ID       int32
	Username string
	Fullname string
	Bio      string
}

func (as AccountService) UpdateByID(ctx context.Context, params AccountUpdateByIDParams) error {
	if len(params.Username) < 3 {
		return ValidationError{
			Field:   "username",
			Message: "shorter than 3 chars",
		}
	} else if len(params.Username) > 50 {
		return ValidationError{
			Field:   "username",
			Message: "longer than 50 chars",
		}
	} else if !reUsername.MatchString(params.Username) {
		return ValidationError{
			Field:   "username",
			Message: "contains chars other than alphanumeric and underscore",
		}
	}
	if len(params.Fullname) > 100 {
		return ValidationError{
			Field:   "fullname",
			Message: "longer than 100 chars",
		}
	}

	err := as.Queries.UpdateAccountDetails(ctx, queries.UpdateAccountDetailsParams{
		ID:       params.ID,
		Username: params.Username,
		Fullname: params.Fullname,
		Bio:      params.Bio,
	})
	if pgErr := (&pgconn.PgError{}); errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrUsernameTaken
	}
	if err != nil {
		return err
	}

	return nil
}

type AccountDeleteParams struct {
	Username string
	Password []byte
}

func (as *AccountService) DeleteByUsername(ctx context.Context, params AccountDeleteParams) error {
	if len(params.Username) < 3 {
		return ValidationError{
			Field:   "username",
			Message: "shorter than 3 chars",
		}
	} else if len(params.Username) > 50 {
		return ValidationError{
			Field:   "username",
			Message: "longer than 50 chars",
		}
	} else if !reUsername.MatchString(params.Username) {
		return ValidationError{
			Field:   "username",
			Message: "contains chars other than alphanumeric and underscore",
		}
	}

	authInfo, err := as.Queries.GetAccountAuthInfo(ctx, params.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrAccountNotFound
	}
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(authInfo.Password, params.Password)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrAccountNotFound
	}
	if err != nil {
		return err
	}

	if err := as.Queries.DeleteAccountByID(ctx, authInfo.ID); err != nil {
		return err
	}

	return nil
}
