package store

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/andrieee44/hackusc/domain"
)

type MemStore struct {
	rw             sync.RWMutex
	userIDTable    map[int64]domain.User
	userEmailTable map[string]int64
	nextID         int64
}

func NewMemStore() *MemStore {
	return &MemStore{
		userIDTable:    make(map[int64]domain.User),
		userEmailTable: make(map[string]int64),
		nextID:         1,
	}
}

func (s *MemStore) CreateUser(_ context.Context, u domain.User) (int64, error) {
	var ok bool

	s.rw.Lock()
	defer s.rw.Unlock()

	_, ok = s.userEmailTable[u.Email]
	if ok {
		return 0, errors.New("email already registered")
	}

	u = copyUser(u)
	u.ID = s.nextID

	s.userIDTable[u.ID] = u
	s.userEmailTable[u.Email] = u.ID
	s.nextID++

	return u.ID, nil
}

func (s *MemStore) GetUserByID(_ context.Context, id int64) (domain.User, error) {
	var (
		u     domain.User
		expID int64
		ok    bool
	)

	s.rw.RLock()
	defer s.rw.RUnlock()

	u, ok = s.userIDTable[id]
	if !ok {
		return domain.User{}, fmt.Errorf("id %d does not exist in database", id)
	}

	expID, ok = s.userEmailTable[u.Email]
	if !ok || expID != id {
		panic("id and email tables not synced")
	}

	return copyUser(u), nil
}

func (s *MemStore) GetUserByEmail(_ context.Context, email string) (domain.User, error) {
	var (
		id int64
		u  domain.User
		ok bool
	)

	s.rw.RLock()
	defer s.rw.RUnlock()

	id, ok = s.userEmailTable[email]
	if !ok {
		return domain.User{}, fmt.Errorf("email %q does not exist in database", email)
	}

	u, ok = s.userIDTable[id]
	if !ok {
		panic("id and email tables not synced")
	}

	return copyUser(u), nil
}

func (s *MemStore) UpdateUser(_ context.Context, newUser domain.User) error {
	var (
		oldUser domain.User
		ok      bool
	)

	s.rw.Lock()
	defer s.rw.Unlock()

	oldUser, ok = s.userIDTable[newUser.ID]
	if !ok {
		return fmt.Errorf("id %d does not exist in database", newUser.ID)
	}

	if oldUser.Email != newUser.Email {
		_, ok = s.userEmailTable[newUser.Email]
		if ok {
			return errors.New("email already registered")
		}

		delete(s.userEmailTable, oldUser.Email)
		s.userEmailTable[newUser.Email] = newUser.ID
	}

	s.userIDTable[newUser.ID] = copyUser(newUser)

	return nil
}

func copyUser(u domain.User) domain.User {
	var newHash []byte

	newHash = make([]byte, len(u.PasswordHash))
	copy(newHash, u.PasswordHash)
	u.PasswordHash = newHash

	return u
}
