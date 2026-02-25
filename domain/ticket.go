package domain

import (
	"maps"
	"slices"
	"time"
)

type Ticket struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Title         string
	Body          string
	CommentIDs    []uint64
	LabelIDs      map[uint64]struct{}
	AttachmentIDs map[uint64]struct{}
	AssigneeIDs   map[uint64]struct{}
	CommunityIDs  map[uint64]struct{}
	CreatorID     uint64
	StateID       uint64
	ID            uint64
}

func (t Ticket) DeepCopy() Ticket {
	return Ticket{
		t.CreatedAt,
		t.UpdatedAt,
		t.Title,
		t.Body,
		slices.Clone(t.CommentIDs),
		maps.Clone(t.LabelIDs),
		maps.Clone(t.AttachmentIDs),
		maps.Clone(t.AssigneeIDs),
		maps.Clone(t.CommunityIDs),
		t.CreatorID,
		t.StateID,
		t.ID,
	}
}
