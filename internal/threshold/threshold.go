// Package threshold drops or flags log entries whose numeric field values
// fall outside a configured range.
package threshold

import (
	"fmt"
	"strconv"

	"logslice/internal/parser"
)

// Checker tests a named numeric field against min/max bounds.
type Checker struct {
	field  string
	minVal float64
	maxVal float64
	hasMin bool
	hasMax bool
}

// New returns a Checker for the given field. Pass a non-nil minVal or maxVal
// pointer to activate that bound; nil means unbounded on that side.
func New(field string, minVal, maxVal *float64) (*Checker, error) {
	if field == "" {
		return nil, fmt.Errorf("threshold: field name must not be empty")
	}
	c := &Checker{field: field}
	if minVal != nil {
		c.minVal = *minVal
		c.hasMin = true
	}
	if maxVal != nil {
		c.maxVal = *maxVal
		c.hasMax = true
	}
	if c.hasMin && c.hasMax && c.minVal > c.maxVal {
		return nil, fmt.Errorf("threshold: min %.4g > max %.4g", c.minVal, c.maxVal)
	}
	return c, nil
}

// Allow returns true when the entry's field value is within the configured
// bounds, or when the field is absent / non-numeric (pass-through).
func (c *Checker) Allow(e parser.Entry) bool {
	v, ok := e.Fields[c.field]
	if !ok {
		return true
	}
	f, err := toFloat(v)
	if err != nil {
		return true
	}
	if c.hasMin && f < c.minVal {
		return false
	}
	if c.hasMax && f > c.maxVal {
		return false
	}
	return true
}

func toFloat(v any) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	}
	return 0, fmt.Errorf("not numeric")
}
