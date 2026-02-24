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
