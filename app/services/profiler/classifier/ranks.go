package classifier

// Ranks contains a list of ranks
type Ranks []Rank

// Sort by Value, Weight
func (r Ranks) Len() int { return len(r) }
func (r Ranks) Less(i, j int) bool {
	if r[i].Value == r[j].Value {
		return r[i].Weight < r[j].Weight
	}
	return r[i].Value < r[j].Value
}
func (r Ranks) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

// Rank contains information about a rank
type Rank struct {
	Name   string
	Value  float64
	Weight int
}

// NewRank creates a new rank entry
func NewRank(name string, value float64, weight int) Rank {
	return Rank{
		Name:   name,
		Value:  value,
		Weight: weight,
	}
}
