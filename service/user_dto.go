package service

import "github.com/andrieee44/hackusc/domain"

type UserDTO struct {
	Name  string
	Email string
	Roles []string
	ID    int64
}

func newUserDTO(u domain.User) UserDTO {
	return UserDTO{u.Name, u.Email, u.Roles, u.ID}
}
