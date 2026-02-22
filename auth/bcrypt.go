package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct{}

func (Bcrypt) Hash(plain string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
}

func (Bcrypt) ComparePlainToHash(plain string, hash []byte) (bool, error) {
	var err error

	err = bcrypt.CompareHashAndPassword(hash, []byte(plain))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
