package domain

import "github.com/google/uuid"

type Role struct {
	Name        string
	Description string
	Permissions uint64
	ID          uuid.UUID
}

func (r Role) DeepCopy() Role {
	return Role{
		r.Name,
		r.Description,
		r.Permissions,
		r.ID,
	}
}
