package comparator

import (
	"encoding/json"
	"reflect"
)

// Equal comparer
type Equal struct{}

// Greater comparer
type Greater struct {
	Equal bool `json:"equal"`
}

// Lesser Comparer
type Lesser struct {
	Equal bool `json:"equal"`
}

// Compare equal imp
func (e *Equal) Compare(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// Compare greater imp
func (g *Greater) Compare(a, b interface{}) bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	switch a.(type) {
	case int:
		aInt, _ := a.(int)
		bInt, _ := b.(int)
		if g.Equal {
			return aInt >= bInt
		}
		return aInt > bInt
	case int64:
		aInt64, _ := a.(int64)
		bInt64, _ := b.(int64)
		if g.Equal {
			return aInt64 >= bInt64
		}
		return aInt64 > bInt64
	case float64:
		aFloat64, _ := a.(float64)
		bFloat64, _ := b.(float64)
		if g.Equal {
			return aFloat64 >= bFloat64
		}
		return aFloat64 > bFloat64
	default:
		return false
	}

}

// Compare lesser imp
func (l *Lesser) Compare(a, b interface{}) bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	switch a.(type) {
	case int:
		aInt, _ := a.(int)
		bInt, _ := b.(int)
		if l.Equal {
			return aInt <= bInt
		}
		return aInt < bInt
	case int64:
		aInt64, _ := a.(int64)
		bInt64, _ := b.(int64)
		if l.Equal {
			return aInt64 <= bInt64
		}
		return aInt64 < bInt64
	case float64:
		aFloat64, _ := a.(float64)
		bFloat64, _ := b.(float64)
		if l.Equal {
			return aFloat64 <= bFloat64
		}
		return aFloat64 < bFloat64
	default:
		return false
	}

}

// MarshalJSON ...
func (e *Equal) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Comp string `json:"type"`
	}{
		"eq",
	})
}

// MarshalJSON ...
func (l *Lesser) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Comp  string `json:"type"`
		Equal bool   `json:"equal"`
	}{
		"lt",
		l.Equal,
	})
}

// MarshalJSON ...
func (g *Greater) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Comp  string `json:"type"`
		Equal bool   `json:"equal"`
	}{
		"gt",
		g.Equal,
	})
}
