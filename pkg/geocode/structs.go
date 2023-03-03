package geocode

type ZipResponse struct {
	Zip     string  `json:"zip"`
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
}

type DirectResponse struct {
	Entities []*DirectResponseEntity `json:"entities"`
}

type DirectResponseEntity struct {
	Name       string            `json:"name"`
	Lat        float64           `json:"lat"`
	Lon        float64           `json:"lon"`
	Country    string            `json:"country"`
	LocalNames map[string]string `json:"local_names"`
}
