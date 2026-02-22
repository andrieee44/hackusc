package service

import "github.com/andrieee44/hackusc/domain"

type UserActor struct {
	ID int64
}

func newUserActor(u domain.User) UserActor {
	return UserActor{u.ID}
}
