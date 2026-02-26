package domain

import (
	"maps"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	Address       Address
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Title         string
	Body          string
	CommentIDs    []uuid.UUID
	Location      GeoPoint
	LabelIDs      map[uuid.UUID]struct{}
	AttachmentIDs map[uuid.UUID]struct{}
	AssigneeIDs   map[uuid.UUID]struct{}
	CommunityIDs  map[uuid.UUID]struct{}
	CreatorID     uuid.UUID
	StateID       uuid.UUID
	ID            uuid.UUID
}

func (t Ticket) DeepCopy() Ticket {
	return Ticket{
		t.Address.DeepCopy(),
		t.CreatedAt,
		t.UpdatedAt,
		t.Title,
		t.Body,
		slices.Clone(t.CommentIDs),
		t.Location.DeepCopy(),
		maps.Clone(t.LabelIDs),
		maps.Clone(t.AttachmentIDs),
		maps.Clone(t.AssigneeIDs),
		maps.Clone(t.CommunityIDs),
		t.CreatorID,
		t.StateID,
		t.ID,
	}
}
