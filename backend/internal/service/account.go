package service

import (
	"context"
	"errors"
	"garlip/internal/queries"
	"regexp"

	"github.com/jackc/pgx/v5/pgconn"
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
	if len(params.Fullname) > 100 {
		return ValidationError{
			Param: "fullname",
			Msg:   "longer than 100 chars",
		}
	}

	err := as.Queries.UpdateAccount(ctx, queries.UpdateAccountParams{
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
