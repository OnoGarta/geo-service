package address

type SearchRequest struct {
	Query string `json:"query" validate:"required,min=2"`
}

type GeocodeRequest struct {
	Lat string `json:"lat" validate:"required"`
	Lng string `json:"lng" validate:"required"`
}

type Response struct {
	Addresses []*Address `json:"addresses"`
}

type Address struct {
	City   string `json:"city,omitempty"`
	Street string `json:"street,omitempty"`
	House  string `json:"house,omitempty"`
	Lat    string `json:"lat,omitempty"`
	Lon    string `json:"lon,omitempty"`
}
