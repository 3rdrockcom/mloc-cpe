package weekday

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
)

var Debug bool

type Weekday struct {
	t models.Transactions
}

func NewWeekday(t models.Transactions) *Weekday {
	w := &Weekday{
		t: t,
	}

	return w
}

func (w *Weekday) Run() Results {
	return w.calc()
}

func (w *Weekday) calc() Results {
	list := make(Results, 7)

	count := 0
	total := 0.0
	for i := range w.t {
		weekday := int(w.t[i].DateTime.Weekday())

		list[weekday].Weekday = w.t[i].DateTime.Weekday()
		list[weekday].Total += w.t[i].Credit
		list[weekday].Count++

		total += w.t[i].Credit
		count++
	}

	for i := range list {
		entry := list[i]
		entry.Probability = entry.Total / total

		list[i] = entry
	}

	return list
}
