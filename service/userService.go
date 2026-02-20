package service

import (
	"context"

	"github.com/andrieee44/hackusc/domain"
	"github.com/andrieee44/hackusc/store"
)

type UserService struct {
	userStore store.UserStore
}

func NewUserService(userStore store.UserStore) *UserService {
	return &UserService{userStore: userStore}
}

func (s *UserService) CreateUser(ctx context.Context, opts ...domain.UserOpt) (int64, error) {
	var (
		u   domain.User
		err error
	)

	u, err = domain.NewUser(opts...)
	if err != nil {
		return 0, err
	}

	return s.userStore.CreateUser(ctx, u)
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
	return s.userStore.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return s.userStore.GetUserByEmail(ctx, email)
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, opts ...domain.UserOpt) error {
	var (
		u   domain.User
		err error
	)

	u, err = s.userStore.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	err = u.Update(opts...)
	if err != nil {
		return err
	}

	return s.userStore.UpdateUser(ctx, u)
}
