package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/andrieee44/hackusc/domain"
)

type UserStore interface {
	Create(context.Context, domain.User) (int64, error)
	GetByID(context.Context, int64) (domain.User, error)
	GetByEmail(context.Context, string) (domain.User, error)
	Update(context.Context, domain.User) error
}

type UserAuth interface {
	Hash(string) ([]byte, error)
	ComparePlainToHash(string, []byte) (bool, error)
}

type UserService struct {
	userStore UserStore
	userAuth  UserAuth
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserDoesNotExist  = errors.New("user does not exist")
	ErrUserLacksRole     = errors.New("user is lacking required role")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrInternal          = errors.New("internal server error")
)

func NewUserService(userStore UserStore, userAuth UserAuth) *UserService {
	return &UserService{
		userStore: userStore,
		userAuth:  userAuth,
	}
}

func (s *UserService) Create(
	ctx context.Context,
	name, email, password string,
) (UserActor, error) {
	var (
		hash []byte
		u    domain.User
		id   int64
		err  error
	)

	hash, err = s.userAuth.Hash(password)
	if err != nil {
		return UserActor{}, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	u, err = domain.NewUser(name, email, hash)
	if err != nil {
		return UserActor{}, err
	}

	id, err = s.userStore.Create(ctx, u)
	if err != nil {
		return UserActor{}, err
	}

	u.ID = id

	return newUserActor(u), nil
}

func (s *UserService) GetByID(
	ctx context.Context,
	id int64,
) (UserDTO, error) {
	var (
		u   domain.User
		err error
	)

	u, err = s.userStore.GetByID(ctx, id)
	if err != nil {
		return UserDTO{}, err
	}

	return newUserDTO(u), nil
}

func (s *UserService) GetByEmail(
	ctx context.Context,
	email string,
) (UserDTO, error) {
	var (
		u   domain.User
		err error
	)

	u, err = s.userStore.GetByEmail(ctx, email)
	if err != nil {
		return UserDTO{}, err
	}

	return newUserDTO(u), nil
}

func (s *UserService) Update(
	ctx context.Context,
	actor UserActor,
	req UserUpdateRequest,
) error {
	var (
		u   domain.User
		err error
	)

	u, err = s.userStore.GetByID(ctx, actor.ID)
	if err != nil {
		return err
	}

	err = req.apply(s.userAuth, &u)
	if err != nil {
		return err
	}

	return s.userStore.Update(ctx, u)
}

func (s *UserService) Login(
	ctx context.Context,
	email, plainPassword string,
) (UserActor, error) {
	var (
		u   domain.User
		ok  bool
		err error
	)

	u, err = s.userStore.GetByEmail(ctx, email)
	if err != nil {
		return UserActor{}, err
	}

	ok, err = s.userAuth.ComparePlainToHash(plainPassword, u.PasswordHash)
	if err != nil {
		return UserActor{}, err
	}

	if !ok {
		return UserActor{}, ErrIncorrectPassword
	}

	return newUserActor(u), nil
}
