package domain

import (
	"maps"
	"slices"
	"time"
)

type User struct {
	CreatedAt    time.Time
	Name         string
	Email        string
	PasswordHash []byte
	RoleIDs      map[uint64]struct{}
	CommunityIDs map[uint64]struct{}
	AvatarID     uint64
	ID           uint64
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
