package domain

type Role struct {
	Name        string
	Description string
	Permissions uint64
	ID          uint64
}

func (r Role) DeepCopy() Role {
	return Role{
		r.Name,
		r.Description,
		r.Permissions,
		r.ID,
	}
}
