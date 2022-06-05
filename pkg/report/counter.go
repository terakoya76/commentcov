package report

import (
	"github.com/terakoya76/commentcov/proto"
)

// scopedCounter holds counters by scope.
type scopedCounter map[string]*Counter

// Counter is a coverage conuter.
type Counter struct {
	Covered int
	Total   int
}

// NewCounter returns an initialized Counter.
func NewCounter() *Counter {
	return &Counter{
		Covered: 0,
		Total:   0,
	}
}

// CalcRate returns a rate from the current counter's state.
func (c *Counter) CalcRate() float64 {
	var t int
	if c.Total == 0 {
		t = 1
	} else {
		t = c.Total
	}

	return 100.0 * float64(c.Covered) / float64(t)
}

// Profile returns a rate from the current counter's state.
func (c *Counter) Profile(item *proto.CoverageItem) {
	if len(item.HeaderComments) > 0 {
		c.Covered++
		c.Total++
	} else {
		c.Total++
	}
}
