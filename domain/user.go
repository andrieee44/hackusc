package domain

type User struct {
	Name         string
	Email        string
	PasswordHash []byte
	Roles        []string
	ID           int64
}

func NewUser(name string, email string, passwordHash []byte) (User, error) {
	var (
		u   User
		err error
	)

	err = u.ApplyUserOpts(
		WithName(name),
		WithEmail(email),
		WithPasswordHash(passwordHash),
		WithRoles([]string{}),
	)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (u *User) ApplyUserOpts(opts ...UserOpt) error {
	var (
		opt UserOpt
		err error
	)

	for _, opt = range opts {
		err = opt(u)
		if err != nil {
			return err
		}
	}

	return nil
}
