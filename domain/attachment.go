package domain

import "time"

type Attachment struct {
	CreatedAt time.Time
	Filename  string
	MIMEType  string
	URL       string
	CreatorID uint64
	ID        uint64
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
