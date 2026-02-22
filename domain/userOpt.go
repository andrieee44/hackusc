package domain

import (
	"errors"
	"fmt"
	"net/mail"
)

type UserOpt func(*User) error

var ErrInvalidArgument = errors.New("invalid argument")

func WithName(name string) UserOpt {
	return func(u *User) error {
		if name == "" {
			return fmt.Errorf("%w: name is empty", ErrInvalidArgument)
		}

		u.Name = name

		return nil
	}
}

func WithEmail(email string) UserOpt {
	return func(u *User) error {
		var (
			addr *mail.Address
			err  error
		)

		addr, err = mail.ParseAddress(email)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
		}

		u.Email = addr.Address

		return nil
	}
}

func WithPasswordHash(passwordHash []byte) UserOpt {
	return func(u *User) error {
		if len(passwordHash) == 0 {
			return fmt.Errorf("%w: passwordHash is empty", ErrInvalidArgument)
		}

		u.PasswordHash = passwordHash

		return nil
	}
}

func WithRoles(roles []string) UserOpt {
	return func(u *User) error {
		if roles == nil {
			return fmt.Errorf("%w: roles is nil", ErrInvalidArgument)
		}

		u.Roles = roles

		return nil
	}
}
