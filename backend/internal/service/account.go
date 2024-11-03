package service

import (
	"context"
	"garlip/internal/queries"
	"regexp"
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

	if err := as.Queries.UpdateAccount(ctx, queries.UpdateAccountParams{
		ID:       params.ID,
		Username: params.Username,
		Fullname: params.Fullname,
		Bio:      params.Bio,
	}); err != nil {
		return err
	}

	return nil
}
