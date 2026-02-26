package domain

import "github.com/google/uuid"

type Label struct {
	Name        string
	Description string
	Color       string
	CreatorID   uuid.UUID
	ID          uuid.UUID
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
