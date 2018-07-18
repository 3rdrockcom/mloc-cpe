package weekday

import (
	"github.com/epointpayment/mloc-cpe/app/models"

	"github.com/shopspring/decimal"
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
	total := decimal.Zero
	for i := range w.t {
		weekday := int(w.t[i].DateTime.Weekday())

		list[weekday].Weekday = w.t[i].DateTime.Weekday()
		list[weekday].Total = list[weekday].Total.Add(w.t[i].Credit)
		list[weekday].Count++

		total = total.Add(w.t[i].Credit)
		count++
	}

	for i := range list {
		entry := list[i]
		entry.Probability = entry.Total.Div(total).Round(3)

		list[i] = entry
	}

	return list
}
