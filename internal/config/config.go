package config

type ConfigResponse struct {
	Results  []Locations `json:"results"`
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
}

type Config struct {
	Next     *string
	Previous *string
}

type Locations struct {
	Name string `json:"name"`
}
