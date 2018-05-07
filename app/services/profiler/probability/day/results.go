package day

type Results []Result

func (r Results) Len() int           { return len(r) }
func (r Results) Less(i, j int) bool { return r[i].Probability > r[j].Probability }
func (r Results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

type Result struct {
	Day         int     `json:"day"`
	Count       int     `json:"count"`
	Total       float64 `json:"total"`
	Probability float64 `json:"probability"`
}
