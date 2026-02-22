package store

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/andrieee44/hackusc/domain"
	"github.com/andrieee44/hackusc/service"
)

type MemStore struct {
	rw          sync.RWMutex
	userIDTable map[int64]domain.User
	nextUserID  int64
}

func NewMemStore() *MemStore {
	return &MemStore{
		userIDTable: make(map[int64]domain.User),
		nextUserID:  1,
	}
}

func (s *MemStore) Create(
	ctx context.Context,
	u domain.User,
) (int64, error) {
	var (
		err error
	)

	_, err = s.GetByEmail(ctx, u.Email)
	if !errors.Is(err, service.ErrUserDoesNotExist) {
		return 0, fmt.Errorf(
			"%w: email %q",
			service.ErrUserAlreadyExists,
			u.Email,
		)
	}

	s.rw.Lock()
	defer s.rw.Unlock()

	u.ID = s.nextUserID
	s.userIDTable[u.ID] = u
	s.nextUserID++

	return u.ID, nil
}

func (s *MemStore) GetByID(
	_ context.Context,
	id int64,
) (domain.User, error) {
	var (
		u  domain.User
		ok bool
	)

	s.rw.RLock()
	defer s.rw.RUnlock()

	u, ok = s.userIDTable[id]
	if !ok {
		return domain.User{}, fmt.Errorf(
			"%w: id %d",
			service.ErrUserDoesNotExist,
			id,
		)
	}

	return u, nil
}

func (s *MemStore) GetByEmail(
	_ context.Context,
	email string,
) (domain.User, error) {
	var u domain.User

	s.rw.RLock()
	defer s.rw.RUnlock()

	for _, u = range s.userIDTable {
		if u.Email == email {
			return u, nil
		}
	}

	return domain.User{}, fmt.Errorf(
		"%w: email %q",
		service.ErrUserDoesNotExist,
		u.Email,
	)
}

func (s *MemStore) Update(
	ctx context.Context,
	newUser domain.User,
) error {
	var (
		oldUser domain.User
		err     error
	)

	oldUser, err = s.GetByID(ctx, newUser.ID)
	if err != nil {
		return err
	}

	if oldUser.Email != newUser.Email {
		_, err = s.GetByEmail(ctx, newUser.Email)
		if err == nil {
			return fmt.Errorf(
				"%w: email %q",
				service.ErrUserAlreadyExists,
				newUser.Email,
			)
		}

		if !errors.Is(err, service.ErrUserDoesNotExist) {
			return err
		}
	}

	s.rw.Lock()
	defer s.rw.Unlock()

	s.userIDTable[newUser.ID] = newUser

	return nil
}
