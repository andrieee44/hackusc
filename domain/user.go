package domain

import (
	"errors"
	"net/mail"
)

type User struct {
	Name         string
	Email        string
	PasswordHash []byte
	ID           int64
}

func NewUser(opts ...UserOpt) (User, error) {
	var (
		u   User
		err error
	)

	ApplyUserOpts(&u, opts...)

	err = u.Validate()
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (u *User) Update(opts ...UserOpt) error {
	ApplyUserOpts(u, opts...)

	return u.Validate()
}

func (u *User) Validate() error {
	var (
		addr *mail.Address
		err  error
	)

	if u.Name == "" {
		return errors.New("name must not be empty")
	}

	addr, err = mail.ParseAddress(u.Email)
	if err != nil {
		return err
	}

	if addr.Address != u.Email {
		return errors.New("email must be a valid address")
	}

	if len(u.PasswordHash) == 0 {
		return errors.New("hash must not be empty")
	}

	return nil
}
