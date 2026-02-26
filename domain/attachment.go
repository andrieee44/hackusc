package domain

import (
	"time"

	"github.com/google/uuid"
)

type Attachment struct {
	CreatedAt time.Time
	Filename  string
	MIMEType  string
	URL       string
	CreatorID uuid.UUID
	ID        uuid.UUID
}

func (a Attachment) DeepCopy() Attachment {
	return Attachment{
		a.CreatedAt,
		a.Filename,
		a.MIMEType,
		a.URL,
		a.CreatorID,
		a.ID,
	}
}
