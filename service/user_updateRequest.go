package service

import (
	"errors"
	"fmt"

	"github.com/andrieee44/hackusc/domain"
)

type UserUpdateRequest struct {
	Name     *string
	Email    *string
	Password *string
}

var ErrEmptyUserUpdateRequest = errors.New("user update request is empty")

func (req UserUpdateRequest) apply(userAuth UserAuth, u *domain.User) error {
	var (
		opts []domain.UserOpt
		hash []byte
		err  error
	)

	if req.Name != nil {
		opts = append(opts, domain.WithName(*req.Name))
	}

	if req.Email != nil {
		opts = append(opts, domain.WithEmail(*req.Email))
	}

	if req.Password != nil {
		hash, err = userAuth.Hash(*req.Password)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInternal, err)
		}

		opts = append(opts, domain.WithPasswordHash(hash))
	}

	if len(opts) == 0 {
		return ErrEmptyUserUpdateRequest
	}

	return u.ApplyUserOpts(opts...)
}
