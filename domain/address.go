package domain

type Address struct {
	Street     string
	Locality   string
	City       string
	Province   string
	PostalCode string
	Country    string
}

func (a Address) DeepCopy() Address {
	return Address{
		a.Street,
		a.Locality,
		a.City,
		a.Province,
		a.PostalCode,
		a.Country,
	}
}
