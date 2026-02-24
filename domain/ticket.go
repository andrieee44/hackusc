package domain

import "time"

type Ticket struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Title         string
	Body          string
	LabelIDs      map[uint64]struct{}
	AttachmentIDs map[uint64]struct{}
	AssigneeIDs   map[uint64]struct{}
	CommunityIDs  map[uint64]struct{}
	CommentIDs    []uint64
	CreatorID     uint64
	StateID       uint64
	ID            uint64
}
