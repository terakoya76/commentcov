package report

import (
	"github.com/commentcov/commentcov/proto"
)

// ScopedCounter holds counters by scope.
type ScopedCounter map[string]*Counter

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

// Merge merges 2 counters into 1.
func (c *Counter) Merge(other *Counter) {
	c.Covered += other.Covered
	c.Total += other.Total
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

// Add adds the values of the given item.
func (c *Counter) Add(item *proto.CoverageItem) {
	if len(item.HeaderComments) > 0 {
		c.Covered++
		c.Total++
	} else {
		c.Total++
	}
}
