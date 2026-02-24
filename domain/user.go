package domain

type User struct {
	Name         string
	Email        string
	PasswordHash []byte
	RoleIDs      map[uint64]struct{}
	AvatarID     uint64
	ID           uint64
}
