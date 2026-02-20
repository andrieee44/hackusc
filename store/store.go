package store

import (
	"context"

	"github.com/andrieee44/hackusc/domain"
)

type UserStore interface {
	CreateUser(context.Context, domain.User) (int64, error)
	GetUserByID(context.Context, int64) (domain.User, error)
	GetUserByEmail(context.Context, string) (domain.User, error)
	UpdateUser(context.Context, domain.User) error
}
