package value

import (
	"encoding/json"
	"errors"
)

// Value ...
type Value struct {
	Value interface{}
	Type  Type
}

// Type ...
type Type string

const (
	// Bool type
	Bool = "bool"
	// Int type
	Int = "int"
	// Int64 type
	Int64 = "int64"
	// Float64 type
	Float64 = "float64"
	// String type
	String = "string"
)

// UnmarshalJSON ...
func (v *Value) UnmarshalJSON(data []byte) error {
	val := struct {
		Type  Type
		Value json.RawMessage
	}{}
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	v.Type = val.Type
	switch val.Type {
	case Bool:
		var b bool
		if err := json.Unmarshal(val.Value, &b); err != nil {
			return err
		}
		v.Value = b
		return nil
	case Int:
		var i int
		if err := json.Unmarshal(val.Value, &i); err != nil {
			return err
		}
		v.Value = i
		return nil
	case Int64:
		var i64 int64
		if err := json.Unmarshal(val.Value, &i64); err != nil {
			return err
		}
		v.Value = i64
		return nil
	case Float64:
		var f64 float64
		if err := json.Unmarshal(val.Value, &f64); err != nil {
			return err
		}
		v.Value = f64
		return nil
	case String:
		var s string
		if err := json.Unmarshal(val.Value, &s); err != nil {
			return err
		}
		v.Value = s
		return nil
	}
	return errors.New("unmarshal failed Value invalid type")
}

// GetValueInterfaces ...
func GetValueInterfaces(values []*Value) []interface{} {
	res := make([]interface{}, len(values))
	for i := range values {
		res[i] = values[i].Value
	}
	return res
}
