package domain

import (
	"maps"
	"slices"
	"time"

	"github.com/google/uuid"
)

type User struct {
	CreatedAt    time.Time
	Name         string
	Email        string
	PasswordHash []byte
	RoleIDs      map[uuid.UUID]struct{}
	CommunityIDs map[uuid.UUID]struct{}
	AvatarID     uuid.UUID
	ID           uuid.UUID
}

func (u User) DeepCopy() User {
	return User{
		u.CreatedAt,
		u.Name,
		u.Email,
		slices.Clone(u.PasswordHash),
		maps.Clone(u.RoleIDs),
		maps.Clone(u.CommunityIDs),
		u.AvatarID,
		u.ID,
	}
}
