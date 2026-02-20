package domain

type UserOpt func(*User)

func ApplyUserOpts(u *User, opts ...UserOpt) {
	var opt UserOpt

	for _, opt = range opts {
		opt(u)
	}
}

func WithName(name string) UserOpt {
	return func(u *User) {
		u.Name = name
	}
}

func WithEmail(email string) UserOpt {
	return func(u *User) {
		u.Email = email
	}
}

func WithPasswordHash(hash []byte) UserOpt {
	return func(u *User) {
		u.PasswordHash = hash
	}
}
