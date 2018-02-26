type HourHist struct {
	Time       int     `json:"time"`
	Close      float64 `json:"close"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Open       float64 `json:"open"`
	Volumefrom float64 `json:"volumefrom"`
	Volumeto   float64 `json:"volumeto"`
}