package domain

type CommunityType struct {
	Name        string
	Description string
	ID          uint64
}

type Community struct {
	Address     Address
	Name        string
	Description string
	Location    GeoPoint
	TypeID      uint64
	ID          uint64
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
