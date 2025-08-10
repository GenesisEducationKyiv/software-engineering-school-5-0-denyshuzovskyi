package weatherapi

type Location struct {
	Name string `json:"name"`
}

type Condition struct {
	Text string `json:"text"`
}

type Current struct {
	LastUpdated int64     `json:"last_updated_epoch"`
	TempC       float32   `json:"temp_c"`
	Condition   Condition `json:"condition"`
	Humidity    int       `json:"humidity"`
}

type CurrentWeather struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}
