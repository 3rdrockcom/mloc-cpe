package profiler

type Results [][]Result

type Result struct {
	Classification string     `json:"classification"`
	Match          float64    `json:"match"`
	Credits        Credits    `json:"credits"`
	Statistics     Statistics `json:"statistics"`
}

type Credits struct {
	AveragePerInterval float64 `json:"average_per_interval"`
	Average            float64 `json:"average"`
}

type Statistics struct {
	Day     []Day     `json:"day"`
	Weekday []Weekday `json:"weekday,omitempty"`
}

type Day struct {
	Day         int     `json:"day"`
	Probability float64 `json:"probability"`
}

type Weekday struct {
	Weekday     string  `json:"weekday"`
	Probability float64 `json:"probability"`
}
