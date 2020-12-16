package kafka

// Reading is a fake temperature reading along with the time it was generated.
type Reading struct {
	Celsius float64 `json:"celsius"`
	Epoch   int64   `json:"epoch"`
}
