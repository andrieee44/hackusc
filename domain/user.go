package domain

import (
	"maps"
	"slices"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Address      Address
	CreatedAt    time.Time
	FirstName    string
	LastName     string
	Email        string
	PasswordHash []byte
	Location     GeoPoint
	Nickname     *string
	MiddleName   *string
	RoleIDs      map[uuid.UUID]struct{}
	CommunityIDs map[uuid.UUID]struct{}
	AvatarID     uuid.UUID
	ID           uuid.UUID
}

func (u User) DeepCopy() User {
	return User{
		u.Address.DeepCopy(),
		u.CreatedAt,
		u.FirstName,
		u.LastName,
		u.Email,
		slices.Clone(u.PasswordHash),
		u.Location.DeepCopy(),
		u.Nickname,
		u.MiddleName,
		maps.Clone(u.RoleIDs),
		maps.Clone(u.CommunityIDs),
		u.AvatarID,
		u.ID,
	}
}
