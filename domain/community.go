package domain

import "github.com/google/uuid"

type CommunityType struct {
	Name        string
	Description string
	ID          uuid.UUID
}

type Community struct {
	Address     Address
	Name        string
	Description string
	Location    GeoPoint
	TypeID      uuid.UUID
	ID          uuid.UUID
}

func (ct CommunityType) DeepCopy() CommunityType {
	return CommunityType{
		ct.Name,
		ct.Description,
		ct.ID,
	}
}

func (c Community) DeepCopy() Community {
	return Community{
		c.Address,
		c.Name,
		c.Description,
		c.Location,
		c.TypeID,
		c.ID,
	}
}
