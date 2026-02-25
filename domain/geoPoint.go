package domain

type GeoPoint struct {
	Latitude  float64
	Longitude float64
}

func (gp GeoPoint) DeepCopy() GeoPoint {
	return GeoPoint{
		gp.Latitude,
		gp.Longitude,
	}
}
