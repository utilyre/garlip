package service

import (
	"context"
	"errors"
	"garlip/internal/postgres"
	"regexp"
)

type AccountService struct {
	Queries *postgres.Queries
}

type AccountUpdateParams struct {
	ID       int32
	Username string
	Fullname string
	Bio      string
}

func (as AccountService) Update(ctx context.Context, params AccountUpdateParams) error {
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

	if len(params.Fullname) == 0 {
		return ValidationError{
			Param: "fullname",
			Err:   errors.New("required"),
		}
	}
	if len(params.Fullname) > 100 {
		return ValidationError{
			Param: "fullname",
			Err:   errors.New("longer than 100 chars"),
		}
	}

	as.Queries.UpdateAccount(ctx, postgres.UpdateAccountParams{
		ID:       params.ID,
		Username: params.Username,
		Fullname: params.Fullname,
		Bio:      params.Bio,
	})
}
