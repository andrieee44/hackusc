package domain

type Label struct {
	Name        string
	Description string
	Color       string
	CreatorID   uint64
	ID          uint64
}

func (l Label) DeepCopy() Label {
	return Label{
		l.Name,
		l.Description,
		l.Color,
		l.CreatorID,
		l.ID,
	}
}
